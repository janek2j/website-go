package main

import (
	"fmt"
	"net/http"

	"github.com/janek2j/website-go/pkg/handlers"
	"github.com/janek2j/website-go/pkg/render"
)

const portNumber = ":5000" // in ":nnnn" format

func main() {
	handlers.Populate()
	render.LoadTemplates()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/users", handlers.UserHandler)
	http.HandleFunc("/user_form", handlers.User_formHandler)
	http.HandleFunc("/items", handlers.ItemHandler)
	http.HandleFunc("/item_form", handlers.Item_formHandler)
	http.HandleFunc("/about", handlers.AboutHandler)
	http.HandleFunc("/other", handlers.OtherHandler)

	fmt.Printf("Starting application on port %s", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
