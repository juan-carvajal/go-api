package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/juan-carvajal/go-api/pkg/api/service/voucher"
	"github.com/juan-carvajal/go-api/pkg/shared/utils"
	"gorm.io/gorm"
)

type ProductService interface {
	handleGetAllProducts() http.HandlerFunc
	handleGetProductByID() http.HandlerFunc
	RegisterRoutes(router *mux.Router) error
}

type DefaultProductService struct {
	productRepo ProductRepo
	voucherRepo voucher.VoucherRepo
}

func (s *DefaultProductService) handleGetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := ProductQueryParams{}

		params.SearchParams = utils.ParseCommonQueryParams(r)

		queryParams := r.URL.Query()

		voucherId := queryParams.Get("voucher")

		voucher, err := s.voucherRepo.GetVoucherByID(voucherId)

		if err == nil {
			params.Voucher = voucher
		}

		products, err := s.productRepo.GetAllProducts(params)
		if err != nil {
			utils.WriteError(w, "error getting products", "failed to get products", err, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
	}
}

func (s *DefaultProductService) handleGetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		rawId, ok := vars["id"]

		if !ok {
			utils.WriteError(w, "no id found", "could not find product id", errors.New("no product id"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(rawId)
		if err != nil {
			utils.WriteError(w, "failed parsing id", "could not parse product id", err, http.StatusInternalServerError)
			return
		}

		product, err := s.productRepo.GetProductByID(id)

		if err == gorm.ErrRecordNotFound {
			utils.WriteError(w, "no product", "product does not exist", gorm.ErrRecordNotFound, http.StatusNotFound)
			return
		}

		if err != nil {
			utils.WriteError(w, "failed getting product", "could not get product", err, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
}

func (s *DefaultProductService) RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/products/{id}", s.handleGetProductByID()).Methods("GET")
	router.HandleFunc("/products", s.handleGetAllProducts()).Methods("GET")

	return nil
}

func NewDefaultProductService(productRepo ProductRepo, voucherRepo voucher.VoucherRepo) ProductService {
	return &DefaultProductService{productRepo, voucherRepo}
}
