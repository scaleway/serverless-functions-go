package handler

import (
	"encoding/json"
	"net/http"
)

// This handle function comes frome our examples and is not modified at all.
func Handle(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"message": "We're all good",
		"healthy": true,
		"number":  4,
		"headers": r.Header,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the header explicitly depending the returned data
	w.Header().Set("Content-Type", "application/json")

	// Customise status code.
	w.WriteHeader(http.StatusOK)

	// Add content to the response
	_, _ = w.Write(responseBytes)
}
