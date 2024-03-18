package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	Controller "github.com/vadim-shalnev/API_Server_dadata/Controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Создаем контекст с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: Query(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидаем сигнала завершения от операционной системы или отмены контекста
	select {
	case <-stop:
		fmt.Println("Получен сигнал завершения")
	case <-ctx.Done():
		fmt.Println("Завершение по таймауту")
	}

	// Выполняем graceful shutdown сервера
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Ошибка при graceful shutdown:", err)
	}
}
func Query() http.Handler {
	r := chi.NewRouter()
	r.Get("/api/login", Controller.Login)
	r.Post("/api/register", Controller.Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(Controller.AuthMiddleware)
		r.Post("/search", Controller.HandleSearch)
		r.Post("/geocode", Controller.HandleGeocode)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))
	return r
}
