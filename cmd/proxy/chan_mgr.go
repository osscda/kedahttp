package main

import (
	"sync"

	nats "github.com/nats-io/nats.go"
)

type chanMgr struct {
	mut         *sync.RWMutex
	writerChan  chan *nats.Msg
	readerChans []chan *nats.Msg
}

func newChanMgr() *chanMgr {
	return &chanMgr{
		mut:         &sync.RWMutex{},
		writerChan:  make(chan *nats.Msg),
		readerChans: nil,
	}
}

func (c *chanMgr) writer() chan<- *nats.Msg {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.writerChan
}

func (c *chanMgr) newReader() <-chan *nats.Msg {
	c.mut.Lock()
	defer c.mut.Unlock()
	newChan := make(chan *nats.Msg)
	c.readerChans = append(c.readerChans, newChan)
	return newChan
}
