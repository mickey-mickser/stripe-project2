package handler

import (
	"context"
	"encoding/json"
	"github.com/mickey-mickser/stripe-project2"
	"net/http"
	"time"
)

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		//405
		h.newErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	defer r.Body.Close()

	var input api.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		//400
		h.newErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	output, err := h.useCase.UserBalance.GetBalance(ctx, input.Username)
	if err != nil {
		//401
		h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"username": input.Username,
		"balance":  output.Balance,
	}); err != nil {
		//500
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

}

// func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) (float64, error) {
//func (h *Handler) getBalanceM(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		username, err := h.getUser(w, r)
//		if err != nil {
//			h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
//			return
//		}
//
//		balance, err := h.useCase.GetBalance(username)
//		if err != nil {
//			h.newErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		if err := json.NewEncoder(w).Encode(balance); err != nil {
//			h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

//	username, err := h.getUser(w, r)
//	if err != nil {
//		h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
//		return 0, err
//	}
//
//	balance, err := h.useCase.GetBalance(username)
//	if err != nil {
//		h.newErrorResponse(w, http.StatusNotFound, err.Error())
//		return 0, err
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(balance); err != nil {
//		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
//		return 0, err
//	}
//	return balance, err
//}
