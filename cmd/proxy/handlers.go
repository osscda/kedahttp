package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	bolt "github.com/boltdb/bolt"
	nats "github.com/nats-io/nats.go"
)

func newForwardingHandler(
	incrementReq func(),
	scaledUpCh <-chan *nats.Msg,
	db *bolt.DB,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go incrementReq()
		// check once to see if there's a container available
		containerURLStr := getContainerEndpoint(db)
		if containerURLStr == "" {
			log.Printf(
				"Waiting for backend container for %s%s",
				r.URL.Host,
				r.URL.String(),
			)
			msg := <-scaledUpCh
			log.Printf("Handler got scaled up event")
			// don't wait for the URL to be in the DB.
			// just get it right away and let the
			// watcher goroutine put it in the DB
			// asynchronously
			containerURLStr = string(msg.Data)
		}
		log.Printf("using container URL %s", containerURLStr)
		// forward the request
		containerURL, err := url.Parse(containerURLStr)
		if err != nil {
			log.Printf(
				"Error parsing container URL string %s (%s)",
				containerURLStr,
				err,
			)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(containerURL)
		proxy.Director = func(req *http.Request) {
			req.URL = containerURL
			req.Host = containerURL.Host
			// req.URL.Scheme = "https"
			// req.URL.Path = r.URL.Path
			// req.URL.Host = containerURL.Host
			// req.URL.Path = containerURL.Path
			reqBytes, _ := httputil.DumpRequest(req, false)
			log.Printf("Proxying request %v", string(reqBytes))
		}
		proxy.ModifyResponse = func(res *http.Response) error {
			respBody, _ := httputil.DumpResponse(res, true)
			log.Printf("Proxied response: %v", string(respBody))
			return nil
		}

		log.Printf("Proxying request to %s to host %s", r.URL.Path, containerURLStr)
		proxy.ServeHTTP(w, r)
	}
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
