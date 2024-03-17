package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api/login", Login)
	r.Post("/api/register", Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/search", Handle)
		r.Post("/geocode", Handle)
	})
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	http.ListenAndServe(":8080", r)
}
