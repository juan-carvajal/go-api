package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/juan-carvajal/go-api/pkg/api/middleware"
	"github.com/juan-carvajal/go-api/pkg/api/service/products"
	"github.com/juan-carvajal/go-api/pkg/api/service/subscriptions"
	"github.com/juan-carvajal/go-api/pkg/api/service/user"
	"github.com/juan-carvajal/go-api/pkg/api/service/voucher"
	"github.com/juan-carvajal/go-api/pkg/migrations"
)

func main() {
	env, err := godotenv.Read()

	if err != nil {
		panic("failed to read .env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env["DB_HOST"], env["DB_USER"], env["DB_PASSWORD"], env["DB_NAME"], env["DB_PORT"])
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("could not init db")
	}

	migrations.AutoMigrateAndSeed(db)

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	productRepo := products.NewDefaultProductRepo(db)
	voucherRepo := voucher.NewDefaultVoucherRepo(db)
	userRepo := user.NewDefaultUserRepo(db)
	subRepo := subscriptions.NewDefaultSubscriptionsRepo(db)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	authRouter := router.PathPrefix("/auth").Subrouter()

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// product routes
	productService := products.NewDefaultProductService(productRepo, voucherRepo)
	productService.RegisterRoutes(apiRouter)

	// users/auth routes
	userService := user.NewDefaultUserService(userRepo)
	userService.RegisterRoutes(authRouter)

	// subscription routes
	subService := subscriptions.NewDefaultSubscriptionService(subRepo, voucherRepo, productRepo)
	subService.RegisterRoutes(apiRouter)

	fmt.Println("starting server")

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
