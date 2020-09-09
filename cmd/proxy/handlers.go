package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newForwardingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svcName := r.URL.Query().Get("name")
		if svcName == "" {
			log.Printf("No service name given")
			w.WriteHeader(400)
			return
		}
		hostPortStr := fmt.Sprintf("http://%s:8080", svcName)
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
			req.URL.Path = r.URL.Path
			req.URL.RawQuery = r.URL.RawQuery
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

func newHealthCheckHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handle Azure Front Door health checks
		w.WriteHeader(http.StatusOK)
	})
}
