package main

import (
	"encoding/json"
	"net/http"
)

// Response structure for API responses.
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Error response for failed requests.
func JSONError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	resp := APIResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}

// Success response for successful requests.
func JSONSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	resp := APIResponse{
		Status: "success",
		Data:   data,
	}
	json.NewEncoder(w).Encode(resp)
}

func JSONFeed(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := APIResponse{
		Status: "success",
		Data:   data,
	}

	json.NewEncoder(w).Encode(resp)
}
