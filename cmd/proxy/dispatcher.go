package main

import (
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

// listens for scaled events from the controller and
// forwards them to the appropriate writer or reader channel
// based on the subject
func startDispatcher(
	nc *nats.Conn,
	scaledUpCh chan<- *nats.Msg,
	scaledDownCh chan<- *nats.Msg,
) {
	timeout := 1 * time.Second
	nc.Subscribe("scaled.*", func(m *nats.Msg) {
		var sendCh chan<- *nats.Msg
		switch m.Subject {
		case "scaled.up":
			sendCh = scaledUpCh
		case "scaled.down":
			sendCh = scaledDownCh
		default:
			log.Printf(
				"The dispatcher doesn't know what to do with subject %s",
				m.Subject,
			)
			return
		}
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case sendCh <- m:
		case <-timer.C:
			log.Printf(
				"Could not dispatch %s message after %s",
				m.Subject,
				timeout,
			)
			return
		}
	})
}
