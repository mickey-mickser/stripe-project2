package handler

import (
	"context"
	"encoding/json"
	"github.com/mickey-mickser/stripe-project2"
	"net/http"
	"time"
)

type gInput struct {
	Username string `json:"username" binding:"required"`
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.newErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var input api.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		//400
		h.newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.useCase.User.CreateUser(ctx, input)
	if err != nil {
		h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})

}

func (h *Handler) getUser /*Username*/ (w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	h.newErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
	//	return
	//}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	defer r.Body.Close()

	var input gInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		//400
		h.newErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	output, err := h.useCase.User.GetUser(ctx, input.Username)
	if err != nil {
		h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

}
