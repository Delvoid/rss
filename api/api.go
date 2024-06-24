package api

import (
	"github.com/Delvoid/go_rss/internal/database"
)

type APIConfig struct {
	DB *database.Queries
}

func NewAPIConfig(db *database.Queries) *APIConfig {
	return &APIConfig{
		DB: db,
	}
}
