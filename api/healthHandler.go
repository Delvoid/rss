package api

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
