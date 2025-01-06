package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mickey-mickser/stripe-project2/pkg/usecase"
)

type Handler struct {
	useCase *usecase.UseCase
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	//

	router.Route("/api", func(r chi.Router) {
		{
			r.Post("/userCreate", h.createUser)
			r.Post("/userGet", h.getUser)
			r.Post("/balance", h.getBalance)
		}

	})
	router.Route("/stripe", func(r chi.Router) {
		{
			r.Get("/{username}/{sum}", h.createPaymentSession)
			//r.Get("/{sessionID}/status", h.getSessionStatus)
		}
	})
	return router
}
