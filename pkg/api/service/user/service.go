package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juan-carvajal/go-api/pkg/auth"
	"github.com/juan-carvajal/go-api/pkg/models"
	"github.com/juan-carvajal/go-api/pkg/shared/utils"
)

type UserService interface {
	handleRegisterNewUser() http.HandlerFunc
	handleGetNewToken() http.HandlerFunc
	RegisterRoutes(router *mux.Router) error
}

type DefaultUserService struct {
	userRepo UserRepo
}

func (s *DefaultUserService) handleRegisterNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var inputUser models.RegisterUser

		err := json.NewDecoder(r.Body).Decode(&inputUser)
		if err != nil {
			utils.WriteError(w, "could not decode request body", "could not decode request body", err, http.StatusBadRequest)
			return
		}

		if len(inputUser.Password) < 5 {

			utils.WriteError(w, "password failed validation", "please use a password with 5 or more characters", err, http.StatusBadRequest)
			return

		}

		pHash, err := auth.HashPassword(inputUser.Password)
		if err != nil {
			utils.WriteError(w, "could not hash password", "user creation failed", err, http.StatusInternalServerError)
			return
		}

		newUser := models.User{
			Username:     inputUser.Username,
			PasswordHash: pHash,
		}

		err = s.userRepo.RegisterUser(newUser)
		if err != nil {
			utils.WriteError(w, "failed registering user", "user creation failed", err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *DefaultUserService) handleGetNewToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputUser models.RegisterUser

		err := json.NewDecoder(r.Body).Decode(&inputUser)
		if err != nil {
			utils.WriteError(w, "could not decode request body", "could not decode request body", err, http.StatusBadRequest)
			return
		}

		dbUser, err := s.userRepo.GetUserByUsername(inputUser.Username)
		if err != nil {
			utils.WriteError(w, "user does not exist", "specified user does not exist", err, http.StatusNotFound)
			return
		}

		pCheck := auth.CheckPasswordHash(inputUser.Password, dbUser.PasswordHash)

		if !pCheck {
			utils.WriteError(w, "password mismatch", "unauthorized", err, http.StatusUnauthorized)
			return
		}

		token, err := auth.CreateToken(uint32(dbUser.ID))
		if err != nil {
			utils.WriteError(w, "failed generating token", "failure generating token", err, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(token))
	}
}

func (s *DefaultUserService) RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/register", s.handleRegisterNewUser()).Methods(http.MethodPost)
	router.HandleFunc("/token", s.handleGetNewToken()).Methods(http.MethodPost)

	return nil
}

func NewDefaultUserService(userRepo UserRepo) UserService {
	return &DefaultUserService{userRepo: userRepo}
}
