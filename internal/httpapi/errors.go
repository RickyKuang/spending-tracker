// Package httpapi wires the HTTP router and handlers for the backend API.
package httpapi

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the JSON envelope returned for every non-2xx response.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// WriteError writes a JSON error envelope with the given HTTP status, machine-readable code,
// and human-readable message.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}

// WriteJSON writes v as a JSON response body with the given HTTP status.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
