package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

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

	// Start the worker for fetching feeds
	w := NewWorker(dbQueries, 60*time.Second)
	go w.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", apiCfg.HealthHandler)
	mux.HandleFunc("GET /v1/err", apiCfg.ErrorHandler)
	mux.HandleFunc("POST /v1/users", apiCfg.CreateUserHandler)
	mux.HandleFunc("GET /v1/users", apiCfg.AuthMiddleware(apiCfg.GetUserHandler))

	mux.HandleFunc("POST /v1/feeds", apiCfg.AuthMiddleware(apiCfg.CreateFeedHandler))
	mux.HandleFunc("GET /v1/feeds", apiCfg.GetAllFeedsHandler)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.AuthMiddleware(apiCfg.CreateFeedFollowHandler))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.AuthMiddleware(apiCfg.GetFeedFollowsHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.AuthMiddleware(apiCfg.DeleteFeedFollowHandler))

	mux.HandleFunc("GET /v1/posts", apiCfg.AuthMiddleware(apiCfg.GetPostsHandler))

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
