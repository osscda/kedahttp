package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

// use approximately the same as the k8s HPA algorithm:
// https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#algorithm-details
func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS (%s)", err)
	}

	// the reqcounter subjects come in from the proxy
	nc.Subscribe("reqcounter", func(m *nats.Msg) {
		log.Printf("reqcounter %v", *m)
		// send dummy scale up event so that the proxy can
		// forward to something for now.
		nc.Publish("scaled.up", []byte("https://gifm.dev"))
		log.Printf("sent scale.up")
	})

	log.Printf("Subscribed to 'reqcounter' topic and waiting")
	select {}
}
