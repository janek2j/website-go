package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/models"
	"github.com/janek2j/website-go/internal/render"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "../../templates"

func getRoutes() http.Handler {
	// what can be put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// initialize sesson
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session // store session in configuration

	// create a template cache
	tempCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}
	app.TemplateCache = tempCache
	app.UseCache = false

	// set the config for the template package
	render.NewRenderer(&app)

	repo := NewRepo(&app)
	NewHandlers(repo)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf) // use CSRF
	mux.Use(SessionLoad)

	// defined routes
	mux.Get("/", CfgRepo.Home)
	mux.Get("/book", CfgRepo.Book)

	mux.Get("/search-availability", CfgRepo.SearchAvailability)
	mux.Post("/search-availability", CfgRepo.PostSearchAvailability)
	// mux.Get("/search-availability-json", CfgRepo.AvailabilityJSON)
	mux.Post("/search-availability-json", CfgRepo.AvailabilityJSON)

	mux.Get("/make-reservation", CfgRepo.MakeReservation)
	mux.Post("/make-reservation", CfgRepo.PostMakeReservation)
	mux.Get("/reservation-summary", CfgRepo.ReservationSummary)

	mux.Get("/contact", CfgRepo.Contact)
	mux.Get("/about", CfgRepo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
// because web server doesn't remember session data - let's save it
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTemplateCache creates template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return tmplCache, err
	}

	// range through all files ending with *page.html
	for _, page := range pages {
		name := filepath.Base(page)
		currentTemplate, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return tmplCache, err
		}

		// combines .page template with .layout template together
		if len(matches) > 0 {
			currentTemplate, err = currentTemplate.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				fmt.Println(err)
				return tmplCache, err
			}
		}

		tmplCache[name] = currentTemplate
	}
	return tmplCache, nil
}
