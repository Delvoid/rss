package api

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, "Internal Server Error", http.StatusInternalServerError)
}
