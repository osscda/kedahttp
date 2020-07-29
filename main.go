package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

	scaledCh := make(chan *nats.Msg, 64)
	if _, err := nc.ChanSubscribe("scaled", scaledCh); err != nil {
		log.Fatalf("Couldn't subscribe to 'scaled' channel")
	}

	go func() {
		for {
			msg := <-scaledCh
			err := db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte("containers"))
				if err != nil {
					return fmt.Errorf("Couldn't create 'containers' bucket (%s)", err)
				}
				containerURL := msg.Data
				bucket.Put(containerURL, nil)
				return nil
			})
			if err != nil {
				break
			}
		}
	}()

	e := echo.New()
	e.Any("/*", newHomeHandler(
		func() { nc.Publish("reqCounter", nil) },
		scaledCh,
		db,
	))
	// e := echo.New()

	// admin := e.Group("/admin")

	// admin.POST("")
	e.Start(":8080")
}

func newHomeHandler(incrementReq func(), scaledCh <-chan *nats.Msg, db *bolt.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		httpReq := c.Request()
		// check once to see if there's a container available
		containerURL := getContainerEndpoint(db)
		if containerURL == "" {
			log.Printf(
				"Waiting for backend container for %s%s",
				httpReq.URL.Host,
				httpReq.URL.String())
			msg := <-scaledCh
			// don't wait for the URL to be in the DB.
			// just get it right away and let the
			// watcher goroutine put it in the DB
			// asynchronously
			containerURL = string(msg.Data)
		}
		// forward the request
		http.DefaultClient.Do(httpReq)
		incrementReq()
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
