package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf) // use CSRF
	mux.Use(SessionLoad)

	// defined routes
	mux.Get("/", handlers.CfgRepo.Home)
	mux.Get("/book", handlers.CfgRepo.Book)

	mux.Get("/search-availability", handlers.CfgRepo.SearchAvailability)
	mux.Post("/search-availability", handlers.CfgRepo.PostSearchAvailability)
	// mux.Get("/search-availability-json", handlers.CfgRepo.AvailabilityJSON)
	mux.Post("/search-availability-json", handlers.CfgRepo.AvailabilityJSON)

	mux.Get("/make-reservation", handlers.CfgRepo.MakeReservation)
	mux.Post("/make-reservation", handlers.CfgRepo.PostMakeReservation)
	mux.Get("/reservation-summary", handlers.CfgRepo.ReservationSummary)

	mux.Get("/contact", handlers.CfgRepo.Contact)
	mux.Get("/about", handlers.CfgRepo.About)

	mux.Get("/users", handlers.UserH)
	mux.Get("/items", handlers.ItemH)
	// mux.Get("/user_form", handlers.User_form)
	// mux.Get("/item_form", handlers.Item_form)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
