package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	Controller "github.com/vadim-shalnev/API_Server_dadata/Controller"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api/login", Controller.Login)
	r.Post("/api/register", Controller.Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(Controller.AuthMiddleware)
		r.Post("/search", Controller.HandleSearch)
		r.Post("/geocode", Controller.HandleGeocode)
	})
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
