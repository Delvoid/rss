package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Delvoid/go_rss/internal/database"
	"github.com/Delvoid/go_rss/models"
	"github.com/google/uuid"
)

func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		RespondWithError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user := models.DatabaseUserToUser(dbUser)
	RespondWithJSON(w, user, http.StatusCreated)
}

func (cfg *APIConfig) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	respondUser := models.DatabaseUserToUser(user)
	RespondWithJSON(w, respondUser, http.StatusOK)
}
