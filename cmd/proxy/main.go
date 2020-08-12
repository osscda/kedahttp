package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	clusterID = "test-cluster"
	clientID  = "cscaler-client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()

	redisHost := os.Getenv("CSCALER_REDIS_MASTER_SERVICE_HOST")
	redisPort := os.Getenv("CSCALER_REDIS_MASTER_SERVICE_PORT")
	if redisHost == "" || redisPort == "" {
		log.Fatal("CSCALER_REDIS_MASTER_SERVICE_HOST or CSCALER_REDIS_MASTER_SERVICE_PORT not found")
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisCl := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "YWRtaW4=", // no password set
		DB:       0,          // use default DB
	})

	pingTimeout, done := context.WithTimeout(ctx, 200*time.Millisecond)
	pingStatus := redisCl.Ping(pingTimeout)
	done()

	if pingStatus.Err() != nil {
		log.Fatalf(
			"Couldn't connect to redis (%s)",
			pingStatus.Err(),
		)
	}

	log.Print("Connected to Redis")

	// // process that listens for incoming scale events from the controller
	// // and sends them to the right channel
	// dispatcher := startDispatcher(sc)
	// // process that processes incoming scale events and records the updates
	// // to the internal DB
	// go listener(
	// 	dispatcher.newScaleUpReader(),
	// 	dispatcher.newScaleDownReader(),
	// 	db,
	// )

	mux := http.NewServeMux()
	mux.HandleFunc("/pong", pongHandler)
	mux.HandleFunc("/", newForwardingHandler(
		func() {
			redisCl.RPush(ctx, "scaler")
			log.Printf("pushed to redis list")
		},
		func() {
			redisCl.RPop(ctx, "scaler")
			log.Printf("popped from redis list")
		},
	))

	// admin := e.Group("/admin")

	port := "8080"
	portEnv := os.Getenv("LISTEN_PORT")
	if portEnv != "" {
		port = portEnv
	}
	log.Printf("Listening on port %s", port)
	portStr := fmt.Sprintf(":%s", port)
	// admin.POST("")
	http.ListenAndServe(portStr, mux)
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	log.Printf("/pong incoming request: %v", string(reqBytes))
	w.Write(reqBytes)
}
