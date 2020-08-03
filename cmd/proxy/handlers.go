package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	bolt "github.com/boltdb/bolt"
	"github.com/labstack/echo"
	nats "github.com/nats-io/nats.go"
)

func newForwardingHandler(
	incrementReq func(),
	scaledUpCh <-chan *nats.Msg,
	db *bolt.DB,
) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		go incrementReq()
		httpReq := c.Request()
		// check once to see if there's a container available
		containerURL := getContainerEndpoint(db)
		if containerURL == "" {
			log.Printf(
				"Waiting for backend container for %s%s",
				httpReq.URL.Host,
				httpReq.URL.String(),
			)
			msg := <-scaledUpCh
			log.Printf("Handler got scaled up event")
			// don't wait for the URL to be in the DB.
			// just get it right away and let the
			// watcher goroutine put it in the DB
			// asynchronously
			containerURL = string(msg.Data)
		}
		// forward the request
		http.DefaultClient.Do(httpReq)
		return nil
	})
}

func getContainerEndpoint(db *bolt.DB) string {
	var containerURLs []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("containers"))
		if err != nil {
			return fmt.Errorf("Couldn't create bucket (%s)", err)
		}
		bucket.ForEach(func(k, v []byte) error {
			containerURLs = append(containerURLs, string(k))
			return nil
		})
		return nil
	})

	if err != nil {
		return ""
	}

	// no real load balancing for now. just get a random container
	// endpoint
	idx := rand.Intn(len(containerURLs))
	return containerURLs[idx]
}
