package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jaimemartinez88/stan.webhook"
	log "github.com/sirupsen/logrus"
)

const (
	httpServerReadTimeout  = 3 * time.Second
	httpServerWriteTimeout = 120 * time.Second
)

var (
	version            string
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

	log.Println("Starting stan.webhook service version: " + version)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.NotFoundHandler = http.HandlerFunc(webhook.HandleNotFound)
	v1 := r.PathPrefix("/v1/").Subrouter()
	v1.HandleFunc("/healthcheck", webhook.HandleHealthcheck).Methods(http.MethodGet)

	v1.HandleFunc("/stan-webhook", webhook.HandleStanWebhook).Methods(http.MethodPost)

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
			"response-code":    m.Code,
		}).Infoln("handler response")
	})
}
