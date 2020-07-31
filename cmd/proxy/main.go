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
	"github.com/labstack/echo"
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

	// TODO: I don't think that NATS will broadcast to this channel.
	// multiple goroutines are gonna be waiting on this channel.
	// all of them will need to wake up at once.
	//
	// maybe use a different NATS subscription API?
	scaledCh := make(chan *nats.Msg, 64)
	if _, err := nc.ChanSubscribe("scaled", scaledCh); err != nil {
		log.Fatalf("Couldn't subscribe to 'scaled' channel")
	}

	go listener(nc, db)

	// TODO: listen for scale-down events from the controller

	e := echo.New()
	e.GET("/pong", pongHandler)
	e.Any("/*", newForwardingHandler(
		func() { nc.Publish("reqCounter", nil) },
		scaledCh,
		db,
	))
	// e := echo.New()

	// admin := e.Group("/admin")

	port := "8080"
	portEnv := os.Getenv("LISTEN_PORT")
	if portEnv != "" {
		port = portEnv
	}
	log.Printf("Listening on port %s", port)
	// admin.POST("")
	e.Start(fmt.Sprintf(":%s", port))
}

func pongHandler(c echo.Context) error {
	reqBytes, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		return err
	}
	log.Printf("/pong incoming request: %v", string(reqBytes))
	return c.String(http.StatusOK, string(reqBytes))
}
