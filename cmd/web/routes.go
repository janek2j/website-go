package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janek2j/website-go/pkg/config"
	"github.com/janek2j/website-go/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// defined routes
	mux.Get("/", handlers.CfgRepo.HomeHandler)
	mux.Get("/about", handlers.CfgRepo.AboutHandler)
	mux.Get("/users", handlers.UserHandler)
	mux.Get("/items", handlers.ItemHandler)
	mux.Get("/other", handlers.OtherHandler)
	// mux.Get("/user_form", handlers.User_formHandler)
	// mux.Get("/item_form", handlers.Item_formHandler)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
