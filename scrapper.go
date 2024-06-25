package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Link        string    `xml:"link"`
		PubDate     string    `xml:"pubDate"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(db *database.Queries, feed database.Feed) error {
	log.Printf("Fetching feed: %s", feed.Url)

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}
	defer resp.Body.Close()

	var rssFeed RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return fmt.Errorf("error decoding feed: %v", err)
	}

	log.Printf("Found %d items in feed", len(rssFeed.Channel.Items))

	for _, item := range rssFeed.Channel.Items {
		publishedAt, err := parseTime(item.PubDate)
		if err != nil {
			// Log the error and continue with the next item
			log.Printf("Error parsing time for item %s: %v", item.Title, err)
			continue
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: sql.NullTime{
				Time:  publishedAt,
				Valid: !publishedAt.IsZero(),
			},
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("Post already exists: %s", item.Title)
			} else {
				log.Printf("Error creating post %s: %v", item.Title, err)
			}
			continue
		}
	}

	return nil
}

func parseTime(timeStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
