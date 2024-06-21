package main

import (
	"log"
	"net/http"

	"github.com/Delvoid/go_rss/api"
	"github.com/Delvoid/go_rss/config"
)

type appConfig struct {
	port string
}

func main() {

	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", api.HealthHandler)
	mux.HandleFunc("GET /v1/err", api.ErrorHandler)

	server := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %s\n", appConfig.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
