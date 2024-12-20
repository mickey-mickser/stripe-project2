package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := errorResponse{Message: message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("Failed to write error response: %v", err)
	}
}
