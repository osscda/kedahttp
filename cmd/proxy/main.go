package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/arschles/containerscaler/externalscaler"
	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	clusterID = "test-cluster"
	clientID  = "cscaler-client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	reqCounter := &reqCounter{i: 0, mut: new(sync.RWMutex)}

	e := echo.New()

	mux := http.NewServeMux()
	middleware := func(fn echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) {
			// TODO: need to figure out a way to get the increment
			// to happen before fn(w, r) happens below. otherwise,
			// the counter won't get incremented right away and the actual
			// handler will hang longer than it needs to
			go func() {
				reqCounter.inc()
			}()
			defer func() {
				reqCounter.dec()
			}()
			fn(c)
		}
	}
	// don't put this inside the middleware so we don't print out an incorrect
	// counter
	e.GET("/counter", func(c echo.Context) error {
		fmt.Fprintf(c.Response(), "%d", reqCounter.get())
		return nil
	})

	// Azure front door health check
	e.GET("/pong", pongHandler)

	e.GET("/azfdhealthz", newHealthCheckHandler())
	e.Any("/", newForwardingHandler())

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting k8s config (%s)", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating k8s clientset (%s)", err)
	}

	e.POST("/admin/deploy", newAdminCreateDeploymentHandler(clientset))
	e.DELETE("/admin/deploy", newAdminDeleteDeploymentHandler(clientset))

	port := "8080"
	portEnv := os.Getenv("LISTEN_PORT")
	if portEnv != "" {
		port = portEnv
	}
	go func() {

		log.Printf("HTTP listening on port %s", port)
		portStr := fmt.Sprintf(":%s", port)
		// admin.POST("")
		log.Fatal(http.ListenAndServe(portStr, mux))
	}()
	go func() {
		log.Printf("GRPC listening on port 9090")
		log.Fatal(startGrpcServer(reqCounter))
	}()
	select {}
}

func startGrpcServer(ctr *reqCounter) error {
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	externalscaler.RegisterExternalScalerServer(grpcServer, newImpl(ctr))
	return grpcServer.Serve(lis)
}

func pongHandler(c echo.Context) error {
	reqBytes, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.String(200, string(reqBytes))
}
