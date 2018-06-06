package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	webhook "github.com/jaimemartinez88/stan-webhook"
	log "github.com/sirupsen/logrus"
)

const (
	headerRealIP = "x-real-ip"

	httpServerReadTimeout  = 3 * time.Second
	httpServerWriteTimeout = 120 * time.Second
)

var (
	corsAllowedHeaders = handlers.AllowedHeaders([]string{"*"})
	corsAllowedDomains = handlers.AllowedOrigins([]string{
		"http://localhost:3000",
	})
	corsAllowedMethods = handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions})
)

func main() {
	// flag handling
	defaultLocation := flag.String("default", "", "location to write a default configuration to (this will overwrite an existing file at this location)")
	configLocation := flag.String("config", "", "JSON config file to load")

	flag.Parse()
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.SetLevel(log.DebugLevel)
	if *defaultLocation == "" && *configLocation == "" {
		log.Println("Using default config:")
		data, _ := json.MarshalIndent(config, "", "  ")
		io.Copy(os.Stdout, bytes.NewReader(data))
		fmt.Printf("\n")
	} else if *defaultLocation != "" {
		writeDefaultConfig(*defaultLocation)
		os.Exit(0)
	} else if *configLocation != "" {
		loadConfig(*configLocation)
	}

	log.Println("Starting stan.webhook service")

	s := webhook.NewService(config.Environment, config.Hostname)

	r := mux.NewRouter()
	r.Use(getRealIPMiddleware)
	r.Use(loggingMiddleware)
	r.NotFoundHandler = http.HandlerFunc(s.NotFound)

	r.HandleFunc("/healthcheck", s.Healthcheck)

	handler := handlers.CORS(corsAllowedHeaders, corsAllowedDomains, corsAllowedMethods)(r)

	httpServer := &http.Server{
		Addr:         ":" + config.HTTP.ListenPort,
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
		Handler:      handler,
	}

	if err := httpServer.ListenAndServe(); nil != err {
		log.Fatalln("Failed to start server", err)
	}

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		l := log.WithFields(log.Fields{
			"request-path":   r.RequestURI,
			"request-method": r.Method,
		})
		l.Infoln()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		m := httpsnoop.CaptureMetrics(next, w, r)
		// next.ServeHTTP(w, r)
		l.WithFields(log.Fields{
			"request-duration": m.Duration,
			"request-ip":       r.Header.Get(headerRealIP),
			"response-code":    m.Code,
		}).Infoln("handler response")
	})
}

func getRealIPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ipAddress := req.RemoteAddr
		fwdAddress := req.Header.Get("X-Forwarded-For") // capitalisation doesn't matter
		if fwdAddress != "" {
			// Got X-Forwarded-For
			ipAddress = fwdAddress // If it's a single IP, then awesome!
			// If we got an array... grab the first IP
			ips := strings.Split(fwdAddress, ", ")
			if len(ips) > 1 {
				ipAddress = ips[0]
			}
		} else {
			host, _, err := net.SplitHostPort(ipAddress)
			if err != nil {
				host = "127.0.0.1"
			}
			if host == "::1" {
				host = "127.0.0.1"
			}
			ipAddress = host
		}
		req.Header.Set(headerRealIP, ipAddress)
		h.ServeHTTP(rw, req)
	})
}
