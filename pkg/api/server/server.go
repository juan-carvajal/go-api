package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/juan-carvajal/go-api/pkg/api/middleware"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	router *mux.Router
}

func (s *Server) routes() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.LoggingMiddleware)

}

func (s *Server) Run() {
	s.routes()
	http.ListenAndServe(":8000", handlers.CompressHandler(s.router))
}
