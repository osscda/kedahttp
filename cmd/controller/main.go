package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const refreshInterval = 500 * time.Millisecond

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS (%s)", err)
	}

	reqCounterChan := make(chan *nats.Msg)
	scaler(nc, reqCounterChan, refreshInterval)

	// the reqcounter messages come through NATS from the proxy
	nc.Subscribe("reqcounter", func(m *nats.Msg) {
		log.Printf("reqcounter %v", *m)
		reqCounterChan <- m
		// send dummy scale up event so that the proxy can
		// forward to something for now.

		// nc.Publish("scaled.up", []byte("https://gifm.dev"))
		log.Printf("sent scale.up")
	})

	log.Printf("Subscribed to 'reqcounter' topic and waiting")
	select {}
}
