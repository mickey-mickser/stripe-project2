package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) createPaymentSession(w http.ResponseWriter, r *http.Request) {
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
	stripe.Key = "sk_test_51QUVTBA95xZSbOgWrPuOOLwRrnfxTMfichpnsoMu0NDIODaVfRXg1fVAbBdPFcZeSBJipCP6chp0xOW4OvGs1fEZ00zdjAR9Kv"

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
	//for beauty
	sessionID := s.ID

	err = h.useCase.Session.CreateSession(r.Context(), sessionID, userUsername.Username, string(s.Status), float64(sum)/100)
	if err != nil {
		logrus.Errorf("Failed to save session for username %s with sessionID %s: %v", username, sessionID, err)
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)

	go h.UpdateSessionStatus(sessionID, sum, username)
}

func (h *Handler) UpdateSessionStatus(sessionID string, sum int64, username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	for {
		time.Sleep(5 * time.Second)

		s, err := session.Get(sessionID, nil)
		if err != nil {
			logrus.Errorf("Failed to get session from Stripe: %v", err)
			continue
		}
		if err := h.useCase.Session.UpdateSessionStatus(ctx, sessionID, string(s.Status)); err != nil {
			logrus.Errorf("Failed to update session status in db: %v", err)
			continue
		}

		if s.Status == stripe.CheckoutSessionStatusComplete {

			userBalance, err := h.useCase.UserBalance.GetBalance(ctx, username)
			if err != nil {
				logrus.Errorf("Failed to get user balance: %v", err)
				continue
			}
			balance := float64(sum)/100 + userBalance.Balance
			if _, err := h.useCase.UserBalance.UpdateUserBalance(ctx, balance, username); err != nil {
				logrus.Errorf("Failed to update user balance: %v", err)
				continue
			}

			//if err := h.useCase.Session.UpdateSessionStatus(ctx, sessionID, string(s.Status)); err != nil {
			//	logrus.Errorf("Failed to update session status in database: %v", err)
			//	continue
			//}
			logrus.Infof("User %s balance updated successfully", username)
			return
		}
	}
}

func (h *Handler) getSessionStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	sessionID := chi.URLParam(r, "sessionID")

	status, err := h.useCase.GetStatus(ctx, sessionID)
	if err != nil {
		//500
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200

	if err := json.NewEncoder(w).Encode(status); err != nil {
		//500
		h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

//
//coin := chi.URLParam(r, "coin")
//
//ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//defer cancel()
////ctx := context.Background()
//cryptoCoin := crypto.CryptoCoin{
//	SymbolFrom: coin,
//}
//
//price, err := h.useCase.GetCoin(ctx, cryptoCoin)
//if err != nil {
//	h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
//	return
//}
//w.Header().Set("Content-Type", "application/json")
//w.WriteHeader(http.StatusOK)
//
//response := map[string]float64{
//	"price": price,
//}
//
//if err := json.NewEncoder(w).Encode(response); err != nil {
//	h.newErrorResponse(w, http.StatusInternalServerError, err.Error())
//}

//
// import (
// import (
//
//	"context"
//	"encoding/json"
//	"github.com/go-chi/chi/v5"
//	"github.com/mickey-mickser/stripe-project2/pkg/cash"
//	"github.com/sirupsen/logrus"
//	"github.com/stripe/stripe-go/v74"
//	"github.com/stripe/stripe-go/v74/checkout/session"
//	"net/http"+
//	"strconv"
//	"time"
//
// )
//v
//
// // Функция для проверки статуса сессии
//
//	func (h *Handler) checkSessionStatus(sessionID string, sum int64, username string) {
//		ctx := context.Background()
//		//defer cancel()
//		balance := float64(sum) / 100
//		for {
//			time.Sleep(5 * time.Second)
//
//			s, err := session.Get(sessionID, nil)
//			if err != nil {
//				logrus.Printf("Error retrieving session: %v", err)
//				continue
//			}
//
//			if s.Status == stripe.CheckoutSessionStatusComplete {
//				usd, err := h.useCase.UserBalance.UpdateUserBalance(ctx, balance, username)
//				h.sessionStorage.UpdateStatus(sessionID, string(s.Status))
//
//				if err != nil {
//					logrus.Printf("Error updating user balance: %v", err)
//					continue
//				}
//				logrus.Printf("Balance %s: %v", username, usd)
//				break
//			}
//		}
//
// }
//
//	func getSessionStatus(w http.ResponseWriter, r *http.Request, cash *cash.SessionStorage) {
//		sessionID := chi.URLParam(r, "sessionID")
//
//		info, exists := cash.Get(sessionID)
//		if !exists {
//			http.Error(w, "Session not found", http.StatusNotFound)
//			return
//		}
//
//		// Возвращаем данные в формате JSON
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(info)
//	}
//package handler
//
//import (
//	"context"
//	"encoding/json"
//	"github.com/go-chi/chi/v5"
//	"github.com/sirupsen/logrus"
//	"github.com/stripe/stripe-go/v74"
//	"github.com/stripe/stripe-go/v74/checkout/session"
//	"net/http"
//	"os"
//	"strconv"
//	"time"
//)
//
//func (h *Handler) createPaymentSession(w http.ResponseWriter, r *http.Request) {
//	username := chi.URLParam(r, "username")
//	sumStr := chi.URLParam(r, "sum")
//	sum, err := strconv.ParseInt(sumStr, 10, 64)
//	if err != nil {
//		logrus.Errorf("Invalid sum parameter: %v", err)
//		http.Error(w, "Invalid sum parameter", http.StatusBadRequest)
//		return
//	}
//
//	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
//	successURL := "https://example.com/success"
//	cancelURL := "https://example.com/cancel"
//
//	params := &stripe.CheckoutSessionParams{
//		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
//		Mode:               stripe.String(string("payment")),
//		SuccessURL:         stripe.String(successURL),
//		CancelURL:          stripe.String(cancelURL),
//		Currency:           stripe.String("usd"),
//		LineItems: []*stripe.CheckoutSessionLineItemParams{
//			{
//				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
//					Currency: stripe.String("usd"),
//					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
//						Name: stripe.String("account replenishment"),
//					},
//					UnitAmount: stripe.Int64(sum),
//				},
//				Quantity: stripe.Int64(1),
//			},
//		},
//	}
//
//	s, err := session.New(params)
//	if err != nil {
//		logrus.Errorf("Checkout session creation failed: %v", err)
//		http.Error(w, "Failed to create payment session", http.StatusInternalServerError)
//		return
//	}
//
//	err = h.useCase.PaymentSession.CreateSession(r.Context(), s.ID, username, float64(sum)/100)
//	if err != nil {
//		logrus.Errorf("Failed to save session: %v", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//
//	http.Redirect(w, r, s.URL, http.StatusSeeOther)
//
//	sessionID := s.ID
//	go h.checkSessionStatus(sessionID, sum, username)
//}
//
//func (h *Handler) checkSessionStatus(sessionID string, sum int64, username string) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
//	defer cancel()
//
//	balance := float64(sum) / 100
//	maxRetries := 12
//	retries := 0
//
//	for retries < maxRetries {
//		retries++
//		time.Sleep(5 * time.Second)
//
//		s, err := session.Get(sessionID, nil)
//		if err != nil {
//			logrus.Printf("Error retrieving session: %v", err)
//			continue
//		}
//
//		if s.Status == stripe.CheckoutSessionStatusComplete {
//			//usd, err := h.useCase.UserBalance.UpdateUserBalance(ctx, balance, username)
//			err := h.useCase.PaymentSession.UpdateStatus(ctx, sessionID, "complete")
//			if err != nil {
//				logrus.Errorf("Failed to update session status: %v", err)
//			}
//			h.sessionStorage.UpdateStatus(sessionID, string(s.Status))
//
//			if err != nil {
//				logrus.Printf("Error updating user balance: %v", err)
//				continue
//			}
//			logrus.Printf("Balance %s updated to: %v", username, usd)
//			break
//		}
//	}
//
//	if retries == maxRetries {
//		logrus.Printf("Session %s check failed after max retries", sessionID)
//	}
//}
//
//func (h *Handler) getSessionStatus(w http.ResponseWriter, r *http.Request) {
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
//}
