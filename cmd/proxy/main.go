package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/arschles/containerscaler/externalscaler"
	"github.com/arschles/containerscaler/pkg/k8s"
	"github.com/arschles/containerscaler/pkg/srv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
)

const (
	clusterID = "test-cluster"
	clientID  = "cscaler-client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	scalerAddress := os.Getenv("CSCALER_SCALER_ADDRESS")
	if scalerAddress == "" {
		log.Fatalf("Need CSCALER_SCALER_ADDRESS")
	}
	log.Printf("Using CSCALER_SCALER_ADDRESS %s", scalerAddress)
	reqCounter := &reqCounter{i: 0, mut: new(sync.RWMutex)}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(userAgentHandler())
	countM := countMiddleware(reqCounter)

	e.Any("/*", newForwardingHandler(), countM)

	adminE := echo.New()
	adminE.Use(middleware.Logger())

	clientset, dynCl, err := k8s.NewClientset()
	if err != nil {
		log.Fatal(err)
	}

	adminE.POST("/app", newAdminCreateAppHandler(clientset, dynCl, scalerAddress))
	adminE.DELETE("/app", newAdminDeleteAppHandler(clientset, dynCl))
	adminE.GET("/pong", pongHandler)
	adminE.GET("/counter", func(c echo.Context) error {
		fmt.Fprintf(c.Response(), "%d", reqCounter.get())
		return nil
	})

	go func() {
		port := fmt.Sprintf(":%s", srv.EnvOr("ADMIN_PORT", "8081"))
		log.Printf("admin server listening on port %s", port)
		log.Fatal(adminE.Start(port))
	}()
	go func() {
		port := fmt.Sprintf(":%s", srv.EnvOr("PROXY_PORT", "8080"))
		log.Printf("proxy listening on port %s", port)
		log.Fatal(e.Start(port))
	}()
	go func() {
		port := fmt.Sprintf(":%s", srv.EnvOr("GRPC_PORT", "9090"))
		log.Printf("GRPC listening on port %s", port)
		log.Fatal(startGrpcServer(port, reqCounter))
	}()
	select {}
}

func startGrpcServer(port string, ctr *reqCounter) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", port))
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
