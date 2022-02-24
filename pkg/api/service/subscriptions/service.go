package subscriptions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/juan-carvajal/go-api/pkg/api/middleware"
	"github.com/juan-carvajal/go-api/pkg/api/service/products"
	"github.com/juan-carvajal/go-api/pkg/api/service/voucher"
	"github.com/juan-carvajal/go-api/pkg/shared/utils"
)

type SubscriptionsService interface {
	handleGetAllUserSubscriptions() http.HandlerFunc
	handleGetSubscriptionById() http.HandlerFunc
	handleCancelSubscription() http.HandlerFunc
	handlePauseSubscription() http.HandlerFunc
	handleUnpauseSubscription() http.HandlerFunc
	handlePostSubscription() http.HandlerFunc
	RegisterRoutes(router *mux.Router) error
}

type DefaultSubscriptionsService struct {
	subscritionRepo SubscriptionRepo
	voucherRepo     voucher.VoucherRepo
	productRepo     products.ProductRepo
}

func (s *DefaultSubscriptionsService) handleGetAllUserSubscriptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		subs, err := s.subscritionRepo.GetUserSubscriptions(uint(id))

		if err != nil {
			utils.WriteError(w, "failed getting subscriptions", "failed getting subscriptions", err, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(subs)
	}
}

func (s *DefaultSubscriptionsService) handleGetSubscriptionById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		uid, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		subId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.WriteError(w, "failed id", "failed parsing id", err, http.StatusBadRequest)
			return
		}

		sub, err := s.subscritionRepo.GetSubscriptionById(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed getting subscription", "failed getting subscription", err, http.StatusInternalServerError)
			return
		}

		if sub.UserID != uint(uid) {
			utils.WriteError(w, "subscription does not belong to user", "unauthorized resource", err, http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(sub)
	}
}

func (s *DefaultSubscriptionsService) handleCancelSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		uid, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		subId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.WriteError(w, "failed id", "failed parsing id", err, http.StatusBadRequest)
			return
		}

		sub, err := s.subscritionRepo.GetSubscriptionById(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed getting subscription", "failed getting subscription", err, http.StatusInternalServerError)
			return
		}

		if sub.UserID != uint(uid) {
			utils.WriteError(w, "subscription does not belong to user", "unauthorized resource", err, http.StatusUnauthorized)
			return
		}

		err = s.subscritionRepo.CancelSubscription(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed cancelling subscription", "failed cancelling subscription", err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *DefaultSubscriptionsService) handlePauseSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		uid, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		subId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.WriteError(w, "failed parsing id", "failed parsing id", err, http.StatusBadRequest)
			return
		}

		sub, err := s.subscritionRepo.GetSubscriptionById(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed getting subscription", "failed getting subscription", err, http.StatusInternalServerError)
			return
		}

		if sub.UserID != uint(uid) {
			utils.WriteError(w, "subscription does not belong to user", "unauthorized resource", err, http.StatusUnauthorized)
			return
		}

		err = s.subscritionRepo.PauseSubscription(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed pausing subscription", "failed pausing subscription", err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *DefaultSubscriptionsService) handleUnpauseSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		uid, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		subId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.WriteError(w, "failed id", "failed parsing id", err, http.StatusBadRequest)
			return
		}

		sub, err := s.subscritionRepo.GetSubscriptionById(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed getting subscription", "failed getting subscription", err, http.StatusInternalServerError)
			return
		}

		if sub.UserID != uint(uid) {
			utils.WriteError(w, "subscription does not belong to user", "unauthorized resource", err, http.StatusUnauthorized)
			return
		}

		err = s.subscritionRepo.UnpauseSubscription(uint(subId))
		if err != nil {
			utils.WriteError(w, "failed unpausing subscription", "failed unpausing subscription", err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *DefaultSubscriptionsService) handlePostSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uid, err := middleware.ExtractUserID(r)

		if err != nil {
			utils.WriteError(w, "error getting user id", "failed to get user id", err, http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)

		productId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			utils.WriteError(w, "failed id", "failed parsing id", err, http.StatusBadRequest)
			return
		}

		product, err := s.productRepo.GetProductByID(int(productId))
		if err != nil {
			utils.WriteError(w, "failed getting product data", "failed getting product data", err, http.StatusInternalServerError)
			return
		}

		params := CreateSubscriptionParams{
			UserID:  uint(uid),
			Product: *product,
		}

		queryParams := r.URL.Query()

		voucherId := queryParams.Get("voucher")

		voucher, err := s.voucherRepo.GetVoucherByID(voucherId)

		if err == nil && voucher.ValidThru.After(time.Now()) {
			params.Voucher = voucher
		}

		err = s.subscritionRepo.CreateSubscription(params)
		if err != nil {
			utils.WriteError(w, "failed creating subscription", "failed creating subscription", err, http.StatusInternalServerError)
			return
		}
	}
}

func (s *DefaultSubscriptionsService) RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/subscriptions", s.handleGetAllUserSubscriptions()).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/{id}", s.handleGetSubscriptionById()).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}/subscribe", s.handlePostSubscription()).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions/{id}/cancel", s.handleCancelSubscription()).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions/{id}/pause", s.handlePauseSubscription()).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions/{id}/unpause", s.handleUnpauseSubscription()).Methods(http.MethodPost)

	return nil
}

func NewDefaultSubscriptionService(subscritionRepo SubscriptionRepo, voucherRepo voucher.VoucherRepo, productRepo products.ProductRepo) SubscriptionsService {
	return &DefaultSubscriptionsService{
		subscritionRepo: subscritionRepo,
		voucherRepo:     voucherRepo,
		productRepo:     productRepo,
	}
}
