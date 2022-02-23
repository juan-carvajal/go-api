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
	"github.com/juan-carvajal/go-api/pkg/api/service/voucher"
	"github.com/juan-carvajal/go-api/pkg/models"
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

	err = db.Debug().AutoMigrate(models.Product{}, models.Voucher{}, models.User{}, models.VoucherRedeem{}, models.Subscription{})
	if err != nil {
		panic("migration failed")
	}

	fmt.Println("migrations completed")

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	productRepo := products.NewDefaultProductRepo(db)
	voucherRepo := voucher.NewDefaultVoucherRepo(db)

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	productService := products.NewDefaultProductService(productRepo, voucherRepo)
	productService.RegisterRoutes(apiRouter)

	fmt.Println("starting server")

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
