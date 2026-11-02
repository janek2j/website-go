package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/janek2j/website-go/pkg/config"
	"github.com/janek2j/website-go/pkg/handlers"
	"github.com/janek2j/website-go/pkg/render"
)

const portNumber = ":8080" // in ":nnnn" format

func main() {
	var app config.AppConfig

	// create a template cache
	tempCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}
	app.TemplateCache = tempCache
	app.UseCache = false

	// set the config for the template package
	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	handlers.Populate()

	http.HandleFunc("/", handlers.Repo.HomeHandler)
	http.HandleFunc("/users", handlers.UserHandler)
	// http.HandleFunc("/user_form", handlers.User_formHandler)
	http.HandleFunc("/items", handlers.ItemHandler)
	// http.HandleFunc("/item_form", handlers.Item_formHandler)
	http.HandleFunc("/about", handlers.Repo.AboutHandler)
	// http.HandleFunc("/other", handlers.OtherHandler)

	fmt.Printf("Starting application on port %s", portNumber)
	err = http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
