package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func newForwardingHandler(
	incrementReq func(),
	// scaledUpCh <-chan *nats.Msg,
	// db *bolt.DB,
) http.HandlerFunc {
	host := os.Getenv("CSCALER_SERVICE_HOST")
	port := os.Getenv("CSCALER_SERVICE_PORT")
	return func(w http.ResponseWriter, r *http.Request) {
		go incrementReq()

		hostPortStr := fmt.Sprintf("%s:%s", host, port)
		log.Printf("using container URL %s", hostPortStr)
		// forward the request
		svcURL, err := url.Parse(hostPortStr)
		if err != nil {
			log.Printf(
				"Error parsing container URL string %s (%s)",
				hostPortStr,
				err,
			)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(svcURL)
		proxy.Director = func(req *http.Request) {
			req.URL = svcURL
			req.Host = svcURL.Host
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

		log.Printf(
			"Proxying request to %s to host %s",
			r.URL.Path,
			hostPortStr,
		)
		proxy.ServeHTTP(w, r)
	}
}
