package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/janek2j/website-go/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not *chi.mux, but is %T", v)
	}
}
