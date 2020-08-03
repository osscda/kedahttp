package main

import (
	"log"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
)

type dispatcher struct {
	mut               *sync.RWMutex
	scaledUpWriter    chan *nats.Msg
	scaledDownWriter  chan *nats.Msg
	scaledUpReaders   []chan *nats.Msg
	scaledDownReaders []chan *nats.Msg
}

// listens for scaled events from the controller and
// forwards them to the appropriate writer or reader channel
// based on the subject
func startDispatcher(
	nc *nats.Conn,
) *dispatcher {

	ret := &dispatcher{
		mut:              new(sync.RWMutex),
		scaledUpWriter:   make(chan *nats.Msg),
		scaledDownWriter: make(chan *nats.Msg),
	}

	// subscriber to ferry nats events into a channel
	timeout := 1 * time.Second
	nc.Subscribe("scaled.*", func(m *nats.Msg) {
		var sendCh chan<- *nats.Msg
		switch m.Subject {
		case "scaled.up":
			sendCh = ret.scaledUpWriter
		case "scaled.down":
			sendCh = ret.scaledDownWriter
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

	// up/down listener goroutine
	go func() {
		for {
			select {
			case msg := <-ret.scaledUpWriter:
				ret.mut.RLock()
				sendToAll(msg, ret.scaledUpReaders)
				ret.mut.RUnlock()
			case msg := <-ret.scaledDownWriter:
				ret.mut.RLock()
				sendToAll(msg, ret.scaledDownReaders)
				ret.mut.RUnlock()
			}
		}
	}()

	return ret
}

func (d *dispatcher) newScaleUpReader() <-chan *nats.Msg {
	d.mut.Lock()
	defer d.mut.Unlock()
	ch := make(chan *nats.Msg)
	d.scaledUpReaders = append(d.scaledUpReaders, ch)
	return ch
}

func (d *dispatcher) newScaleDownReader() <-chan *nats.Msg {
	d.mut.Lock()
	defer d.mut.Unlock()
	ch := make(chan *nats.Msg)
	d.scaledDownReaders = append(d.scaledDownReaders, ch)
	return ch
}

func sendToAll(msg *nats.Msg, chans []chan *nats.Msg) {
	for _, ch := range chans {
		go func(ch chan *nats.Msg) {
			t := time.NewTimer(1 * time.Second)
			select {
			case ch <- msg:
			case <-t.C:
				// TODO: log timeout - couldn't send message within 1 second
			}
		}(ch)
	}
}
