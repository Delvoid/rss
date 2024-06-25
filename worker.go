package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Delvoid/go_rss/internal/database"
)

type Worker struct {
	db     *database.Queries
	ticker *time.Ticker
}

func NewWorker(db *database.Queries, interval time.Duration) *Worker {
	return &Worker{
		db:     db,
		ticker: time.NewTicker(interval),
	}
}

func (w *Worker) Start() {
	log.Println("Starting worker...")
	for ; ; <-w.ticker.C {
		w.FetchFeeds()
	}
}

func (w *Worker) FetchFeeds() {
	log.Println("Fetching feeds...")
	feeds, err := w.db.GetNextFeedsToFetch(context.Background(), 10)
	if err != nil {
		log.Printf("Error fetching feeds: %v", err)
		return
	}

	log.Printf("Found %d feeds to fetch", len(feeds))

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()
			err := FetchFeed(w.db, feed)
			if err != nil {
				log.Printf("Error fetching feed %s: %v", feed.Url, err)
				return
			}

			_, err = w.db.MarkFeedFetched(context.Background(), feed.ID)
			if err != nil {
				log.Printf("Error marking feed as fetched: %v", err)
			}
		}(feed)
	}
	wg.Wait()
}
