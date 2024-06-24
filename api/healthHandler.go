package api

import (
	"net/http"
)

func (cfg *APIConfig) HealthHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
