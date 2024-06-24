package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Delvoid/go_rss/api"
	"github.com/Delvoid/go_rss/config"
	"github.com/Delvoid/go_rss/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", appConfig.DBUrl)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	apiCfg := api.NewAPIConfig(dbQueries)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", apiCfg.HealthHandler)
	mux.HandleFunc("GET /v1/err", apiCfg.ErrorHandler)
	mux.HandleFunc("POST /v1/users", apiCfg.CreateUserHandler)

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
