package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/juan-carvajal/go-api/pkg/auth"
)

var (
	ErrNoSession = errors.New("no session found on request")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			uid, err := auth.ExtractTokenID(r)

			if err != nil {
				http.Error(w, fmt.Sprintf("Access denied: %s", err.Error()), http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", uid)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func ExtractUserID(r *http.Request) (uint32, error) {
	userId, ok := r.Context().Value("user_id").(uint32)

	if !ok {
		return 0, ErrNoSession
	}

	return userId, nil
}
