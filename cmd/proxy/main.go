package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	bolt "github.com/boltdb/bolt"
	nats "github.com/nats-io/nats.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS (%s)", err)
	}
	db, err := bolt.Open("cscaler.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Error connecting to boltdb (%s)", err)
	}

	// process that listens for incoming scale events from the controller
	// and sends them to the right channel
	dispatcher := startDispatcher(nc)
	// process that processes incoming scale events and records the updates
	// to the internal DB
	go listener(
		dispatcher.newScaleUpReader(),
		dispatcher.newScaleDownReader(),
		db,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/pong", pongHandler)
	mux.HandleFunc("/", newForwardingHandler(
		func() {
			nc.Publish("reqcounter", nil)
			log.Printf("published reqcounter")
		},
		dispatcher.newScaleUpReader(),
		db,
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
