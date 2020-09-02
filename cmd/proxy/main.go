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

	reqCounter := new(int64)
	*reqCounter = 0

	mux := http.NewServeMux()
	middleware := func(fn http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCounter := atomic.AddInt64(reqCounter, 1)
			atomic.StoreInt64(reqCounter, newCounter)
			defer func() {
				newCounter = newCounter - 1
				atomic.StoreInt64(reqCounter, newCounter)
			}()
			fn(w, r)
		})
	}
	// don't put this inside the middleware so we don't print out an incorrect
	// counter
	mux.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", atomic.LoadInt64(reqCounter))
	})
	// Azure front door health check
	mux.HandleFunc("/pong", middleware(pongHandler))

	mux.HandleFunc("/azfdhealthz", newHealthCheckHandler())
	mux.HandleFunc("/", middleware(newForwardingHandler()))

	mux.HandleFunc("/admin/deploy", newAdminDeployHandler())

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
		log.Fatal(startGrpcServer(reqCounter))
	}()
	select {}
}

func startGrpcServer(ctr *int64) error {
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
	w.Write(reqBytes)
}
