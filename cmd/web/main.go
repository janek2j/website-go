package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/janek2j/website-go/pkg/config"
	"github.com/janek2j/website-go/pkg/handlers"
	"github.com/janek2j/website-go/pkg/render"
)

const portNumber = ":8080" // in ":nnnn" format
var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	DbProcess()

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

	db, err := sql.Open("sqlite", "file:app.db")
	if err != nil {
		log.Fatal("cannot open/connect to database")
	}
	defer db.Close()
	app.Db = db

	// set the config for the template package
	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	handlers.Populate()

	fmt.Printf("Starting application on port %s", portNumber)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
