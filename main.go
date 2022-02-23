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
	"github.com/juan-carvajal/go-api/pkg/models"
	"github.com/juan-carvajal/go-api/pkg/models/shared"
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

	err = db.Debug().AutoMigrate(models.Product{}, models.Voucher{})
	if err != nil {
		panic("migration failed")
	}

	fmt.Println("migrated")

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler

		repo := products.NewDefaultProductRepo(db)

		products, err := repo.GetAllProducts(products.ProductQueryParams{
			SearchParams: shared.SearchParams{Pagination: shared.Pagination{PageSize: 100, Offset: 5}},
		})

		if err != nil {
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode(products)
	})

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
