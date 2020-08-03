package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS (%s)", err)
	}
	nc.Subscribe("reqcounter", func(m *nats.Msg) {
		log.Printf("reqcounter %v", *m)
		nc.Publish("scaled.up", []byte("https://gifm.dev"))
		log.Printf("sent scale.up")
	})

	log.Printf("Subscribed to 'reqcounter' topic and waiting")
	select {}
}
