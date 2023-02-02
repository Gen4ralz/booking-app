package main

import (
	"net/http"

	"github.com/gen4ralz/booking-app/pkg/config"
	"github.com/gen4ralz/booking-app/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/about", handlers.Repo.About)

	return mux
}