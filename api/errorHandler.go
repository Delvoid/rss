package api

import (
	"net/http"
)

func (cfg *APIConfig) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, "Internal Server Error", http.StatusInternalServerError)
}
