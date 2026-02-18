package main

import (
	"flag"
	"net/http"

	"github.com/Foga2H/ya-go-gophermart/internal/auth"
	mainConfig "github.com/Foga2H/ya-go-gophermart/internal/config"
	dbConfig "github.com/Foga2H/ya-go-gophermart/internal/config/db"
	"github.com/Foga2H/ya-go-gophermart/internal/gzip"
	handler "github.com/Foga2H/ya-go-gophermart/internal/handler/auth"
	handler2 "github.com/Foga2H/ya-go-gophermart/internal/handler/order"
	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/middleware"
	"github.com/Foga2H/ya-go-gophermart/internal/storage/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	flag.Parse()

	newLogger := logger.NewLogger()

	newConfig := mainConfig.NewConfig()
	newDBConfig := dbConfig.NewConfig()
	dbConnect := db.NewDatabaseConnection(newDBConfig)

	dbConnection, err := dbConnect.Ping()
	if err != nil {
		newLogger.Fatalf("Error connecting to database: %v", err)
	}

	m := db.NewMigration(dbConnection)
	err = m.Up()
	if err != nil {
		newLogger.Fatalf("Error up migration: %s", err)
	}

	newGzip := gzip.NewGzip()
	authToken := auth.NewToken(newConfig.JWTSecret, newLogger)
	jwt := middleware.NewJwt(authToken, newLogger)

	userRepository := db.NewUserRepository(dbConnection, newLogger)
	orderRepository := db.NewOrderRepository(dbConnection, newLogger)

	router.Use(newLogger.Middleware())
	router.Use(newGzip.Middleware())

	router.Route("/api", func(r chi.Router) {
		r.Post("/user/register", handler.NewRegisterUserHandler(userRepository, authToken, newLogger).ServeHTTP)
		r.Post("/user/login", handler.NewLoginUserHandler(userRepository, authToken, newLogger).ServeHTTP)

		r.Group(func(r chi.Router) {
			r.Use(jwt.Middleware())

			r.Post("/user/orders", handler2.NewCreateOrderHandler(orderRepository, newLogger).ServeHTTP)
			r.Get("/user/orders", handler2.NewGetUserOrdersHandler(orderRepository, newLogger).ServeHTTP)
		})
	})

	err = http.ListenAndServe(newConfig.BaseURL, router)
	if err != nil {
		panic(err)
	}
}
