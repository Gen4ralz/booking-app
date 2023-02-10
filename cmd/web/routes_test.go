package main

import (
	"testing"

	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do not thing
	default:
		t.Errorf("type is not *chi.Mux, type is %T", v)
	}
}