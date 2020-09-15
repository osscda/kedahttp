package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	echo "github.com/labstack/echo/v4"
)

// TODO: use proxy handler: https://echo.labstack.com/middleware/proxy ??
func newForwardingHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		svcName := c.Request().URL.Query().Get("name")
		if svcName == "" {
			log.Printf("No service name given")
			return c.String(400, "No service name given")
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
			return c.String(400, fmt.Sprintf("Error parsing URL string %s (%s)", hostPortStr, err))
		}

		r := c.Request()

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
		w := c.Response()
		proxy.ServeHTTP(w, r)
		return nil
	}
}

func newHealthCheckHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// handle Azure Front Door health checks
		return c.String(200, "OK")
	}
}
