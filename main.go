package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type appConfig struct {
	port string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	appConfig := appConfig{
		port: os.Getenv("PORT"),
	}
	// Set default port if not provided
	if appConfig.port == "" {
		appConfig.port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	server := &http.Server{
		Addr:    ":" + appConfig.port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %s\n", appConfig.port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
