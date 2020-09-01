package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"sync/atomic"
	"time"

	"github.com/arschles/containerscaler/externalscaler"
	"google.golang.org/grpc"
)

const (
	clusterID = "test-cluster"
	clientID  = "cscaler-client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// // process that listens for incoming scale events from the controller
	// // and sends them to the right channel
	// dispatcher := startDispatcher(sc)
	// // process that processes incoming scale events and records the updates
	// // to the internal DB
	// go listener(
	// 	dispatcher.newScaleUpReader(),
	// 	dispatcher.newScaleDownReader(),
	// 	db,
	// )

	reqCounter := int64(0)

	mux := http.NewServeMux()
	mux.HandleFunc("/pong", pongHandler)
	mux.HandleFunc("/", newForwardingHandler(
		func() {
			log.Printf("request start")
			newCounter := atomic.AddInt64(&reqCounter, 1)
			atomic.StoreInt64(&reqCounter, newCounter)
		},
		func() {
			log.Printf("request end")
			newCounter := atomic.AddInt64(&reqCounter, -1)
			atomic.StoreInt64(&reqCounter, newCounter)
		},
	))

	// admin := e.Group("/admin")

	port := "8080"
	portEnv := os.Getenv("LISTEN_PORT")
	if portEnv != "" {
		port = portEnv
	}
	go func() {

		log.Printf("HTTP listening on port %s", port)
		portStr := fmt.Sprintf(":%s", port)
		// admin.POST("")
		log.Fatal(http.ListenAndServe(portStr, mux))
	}()
	go func() {
		log.Printf("GRPC listening on port 9090")
		log.Fatal(startGrpcServer())
	}()
	select {}
}

func startGrpcServer(ctr int64) error {
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	externalscaler.RegisterExternalScalerServer(grpcServer, externalscaler.NewImpl(ctr))
	return grpcServer.Serve(lis)
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	log.Printf("/pong incoming request: %v", string(reqBytes))
	w.Write(reqBytes)
}
