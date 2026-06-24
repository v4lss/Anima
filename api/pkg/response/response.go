// Package response provides helpers for uniform JSON HTTP responses.
package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type envelope map[string]any

// JSON writes a JSON payload with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success wraps data in a {"data": ...} envelope.
func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, envelope{"data": data})
}

// Error wraps a message in a {"error": ...} envelope.
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, envelope{"error": msg})
}

// DecodeJSON decodes the request body into the provided destination.
func DecodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

// DecodeJSONLimit decodes the request body with a size limit.
func DecodeJSONLimit(r *http.Request, dst any, limit int64) error {
	limitedReader := io.LimitReader(r.Body, limit)
	return json.NewDecoder(limitedReader).Decode(dst)
}
