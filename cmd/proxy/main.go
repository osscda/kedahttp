package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "cscaler-client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL("localhost:4222"),
	)
	if err != nil {
		log.Fatalf("Error connecting to NATS (%s)", err)
	}

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

	mux := http.NewServeMux()
	mux.HandleFunc("/pong", pongHandler)
	mux.HandleFunc("/", newForwardingHandler(
		func() {
			sc.Publish("reqcounter", nil)
			log.Printf("published reqcounter")
		},
	))

	// admin := e.Group("/admin")

	port := "8080"
	portEnv := os.Getenv("LISTEN_PORT")
	if portEnv != "" {
		port = portEnv
	}
	log.Printf("Listening on port %s", port)
	portStr := fmt.Sprintf(":%s", port)
	// admin.POST("")
	http.ListenAndServe(portStr, mux)
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
