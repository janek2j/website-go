package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/driver"
	"github.com/janek2j/website-go/internal/handlers"
	"github.com/janek2j/website-go/internal/models"
	"github.com/janek2j/website-go/internal/render"
	"github.com/joho/godotenv"
)

const portNumber = ":8080" // in ":nnnn" format
var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// err = db.SQL.Ping()
	// if err != nil {
	// 	log
	// }

	fmt.Printf("Starting application on port %s", portNumber)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// load secret values
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBPASSWORD := os.Getenv("DBPASSWORD")

	// what can be put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Book{})
	gob.Register(models.Borrowing{})

	// change this to true when in production
	app.InProduction = false

	// initialize session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session // store session in configuration

	// connect to database
	log.Println("Connecting to database...")

	connection_string := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=pruntu", DBPASSWORD)
	db, err := driver.ConnectSQL(connection_string)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database")

	// create a template cache
	tempCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
		return nil, err
	}
	app.TemplateCache = tempCache
	app.UseCache = false

	// set the config for the template package
	render.NewRenderer(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}
