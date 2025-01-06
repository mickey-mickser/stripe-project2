package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DataSession struct {
	SessionID string
	Sum       float64
	Username  string
	Status    string
}

func (h *Handler) createPaymentSession(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error write .env: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	username := chi.URLParam(r, "username")
	if username == "" {
		//400
		h.newErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}
	userUsername, err := h.useCase.User.GetUser(ctx, username)
	if err != nil {
		//401
		h.newErrorResponse(w, http.StatusUnauthorized, err.Error())
		logrus.Infof("Username not found: %s", username)
		return
	}
	sum, err := strconv.ParseInt(chi.URLParam(r, "sum"), 10, 64)
	if err != nil {
		//400
		h.newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	stripe.Key = os.Getenv("KEY_STRIPE")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string("payment")),
		SuccessURL:         stripe.String("https://example.com/success"),
		CancelURL:          stripe.String("https://example.com/cancel"),
		Currency:           stripe.String("usd"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("account replenishment"),
					},
					UnitAmount: stripe.Int64(sum),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		logrus.Errorf("Checkout session creation failed: %v", err)
		//500
		http.Error(w, "Checkout session creation failed", http.StatusInternalServerError)
		return
	}
	sessionID := s.ID
	var ds *DataSession
	ds = &DataSession{}
	ds.SessionID = sessionID
	ds.Sum = float64(sum) / 100
	ds.Username = username

	fmt.Println(ds, "111asd")
	err = h.useCase.Session.CreateSession(r.Context(), sessionID, userUsername.Username, string(s.Status), float64(sum)/100)
	if err != nil {
		logrus.Errorf("Failed to save session for username %s with sessionID %s: %v", username, sessionID, err)
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}
