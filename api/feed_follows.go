package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/Delvoid/go_rss/models"
	"github.com/google/uuid"
)

func (cfg *APIConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		RespondWithError(w, "Failed to create feed follow", http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, models.DatabaseFeedFollowToFeedFollow(feedFollow), http.StatusCreated)
}

func (cfg *APIConfig) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, "Failed to get feed follows", http.StatusInternalServerError)
		return
	}

	respondFeedFollows := make([]models.FeedFollow, len(feedFollows))
	for i, dbFeedFollow := range feedFollows {
		respondFeedFollows[i] = models.DatabaseFeedFollowToFeedFollow(dbFeedFollow)
	}

	RespondWithJSON(w, respondFeedFollows, http.StatusOK)
}

func (cfg *APIConfig) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	feedFollowIDStr := r.PathValue("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		RespondWithError(w, "Invalid feed follow ID", http.StatusBadRequest)
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, "Failed to delete feed follow", http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
