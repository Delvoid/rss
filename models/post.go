package models

import (
	"time"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func DatabasePostToPost(dbPost database.Post) Post {
	post := Post{
		ID:        dbPost.ID,
		CreatedAt: dbPost.CreatedAt,
		UpdatedAt: dbPost.UpdatedAt,
		Title:     dbPost.Title,
		Url:       dbPost.Url,
		FeedID:    dbPost.FeedID,
	}

	if dbPost.Description.Valid {
		post.Description = &dbPost.Description.String
	}

	if dbPost.PublishedAt.Valid {
		post.PublishedAt = &dbPost.PublishedAt.Time
	}

	return post
}
