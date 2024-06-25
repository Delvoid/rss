package api

import (
	"net/http"
	"strconv"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/Delvoid/go_rss/models"
)

func (cfg *APIConfig) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	dbPosts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		RespondWithError(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	posts := make([]models.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = models.DatabasePostToPost(dbPost)
	}

	RespondWithJSON(w, posts, http.StatusOK)
}
