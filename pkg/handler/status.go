package handler

//
//import (
//	"encoding/json"
//	"github.com/go-chi/chi/v5"
//	"net/http"
//)
//
//func (h *Handler) getStatus(w http.ResponseWriter, r *http.Request) {
//	sessionID := chi.URLParam(r, "sessionID")
//
//	info, exists := h.sessionStorage.Get(sessionID)
//	if !exists {
//		http.Error(w, "Session not found", http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(info)
//
//}
