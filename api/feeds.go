package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/Delvoid/go_rss/models"
	"github.com/google/uuid"
)

func (cfg *APIConfig) CreateFeedHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		RespondWithError(w, "Failed to create feed", http.StatusInternalServerError)
		return
	}
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		RespondWithError(w, "Failed to create feed follow", http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, map[string]interface{}{
		"feed":        models.DatabaseFeedToFeed(feed),
		"feed_follow": models.DatabaseFeedFollowToFeedFollow(feedFollow),
	}, http.StatusCreated)
}

func (cfg *APIConfig) GetAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		RespondWithError(w, "Failed to get feeds", http.StatusInternalServerError)
		return
	}

	respondFeeds := make([]models.Feed, len(feeds))
	for i, dbFeed := range feeds {
		respondFeeds[i] = models.DatabaseFeedToFeed(dbFeed)
	}

	RespondWithJSON(w, respondFeeds, http.StatusOK)
}
