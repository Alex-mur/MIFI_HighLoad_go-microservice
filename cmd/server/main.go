package main

import (
	"context"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-microservice/internal/handlers"
	"go-microservice/internal/middleware"
	"go-microservice/internal/services"
)

func main() {
	// Инициализация логгера (упрощенный вариант)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Инициализация сервисов
	userService := services.NewUserService()
	auditService := services.NewAuditService()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService, auditService)

	// Создание роутера
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.RateLimitMiddleware)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// Маршруты API
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Маршруты для мониторинга
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.Handle("/metrics", middleware.PrometheusHandler())

	// Настройка сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Could not gracefully shutdown the server")
		}
		close(done)
	}()

	log.Info().Str("port", "8080").Msg("Server is starting")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Could not start server")
	}

	<-done
	log.Info().Msg("Server stopped")
}
