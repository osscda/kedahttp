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

	scaledUp := newChanMgr()
	scaledDown := newChanMgr()
	// process that listens for incoming scale events from the controller
	// and sends them to the right channel
	go startDispatcher(nc, scaledUp.writer(), scaledDown.writer())
	// process that processes incoming scale events and records the updates
	// to the internal DB
	go listener(nc, scaledUp, scaledDown, db)

	e := echo.New()
	e.GET("/pong", pongHandler)
	e.Any("/*", newForwardingHandler(
		func() {
			nc.Publish("reqcounter", nil)
			log.Printf("published reqcounter")
		},
		scaledUp.newReader(),
		db,
	))

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
