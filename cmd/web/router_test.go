package main

import (
	"testing"

	"github.com/dhanekom/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig
	handler := routes(app)

	switch v := handler.(type) {
	case *chi.Mux:
		// Success
	default:
		t.Errorf("expected *chi.Mux, got %T", v)
	}
}
