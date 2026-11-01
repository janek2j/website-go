package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const portNumber = ":5000" // in ":nnnn" format

type User struct {
	Name  string
	Email string
}

type Item struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	IsForRent   bool    `json:"rentable"`
}

type formDataForUser struct {
	*User
	Errors []string
}

type formDataForUserForItem struct {
	*Item
	Errors []string
}

type homeData struct {
	PresentTime   string
	PresentRandom int
}

type AppUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = make([]*User, 0, 10)
var items = make([]*Item, 0, 10)

var welcome = homeData{time.Now().String(), rand.Intn(6) + 1}

var templates = make(map[string]*template.Template, 3)

// Populate fills User and Item collections with some initial data
func Populate() {
	users = append(users,
		&User{"Janek", "janek@example.com"},
		&User{"Ola", "ola@example.com"})

	jsonForItems := `
	[
		{
			"name": "Podstawy elektroniki",
			"category": "books",
			"description": "more notes",
			"price": "30.0",
			"rentable": true
		},
		{
			"name": "Hobbit",
			"category": "books",
			"description": "more notes",
			"price": "65.0",
			"rentable": true
		},
		{
			"name": "Wiertarka Skill",
			"category": "tools",
			"description": "moc 850W",
			"price": "300.0",
			"rentable": false
		}
	]`

	var unmarshalledItems []Item
	err := json.Unmarshal([]byte(jsonForItems), &unmarshalledItems)
	if err != nil {
		log.Println("Error unmarshalling item json", err)
	}
	log.Printf("unmarshalled: %v", unmarshalledItems)
	for _, c := range unmarshalledItems {
		items = append(items, &c)
	}

}

func loadTemplates() {
	templateNames := [5]string{"welcome", "user_list", "item_list", "user_form", "item_form"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			log.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, welcome)
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	templates["user_list"].Execute(writer, users)
}

func itemHandler(writer http.ResponseWriter, request *http.Request) {
	templates["item_list"].Execute(writer, items)
}

func item_formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["item_form"].Execute(writer, formDataForUserForItem{
			Item: &Item{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()

		parsed_price, _ := strconv.ParseFloat(request.Form["price"][0], 32)
		responseData := Item{
			Name:        request.Form["name"][0],
			Category:    request.Form["category"][0],
			Description: request.Form["description"][0],
			Price:       float32(parsed_price),
			IsForRent:   request.Form["isForRent"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please enter item name")
		}
		if responseData.Description == "" {
			errors = append(errors, "Please enter some description")
		}
		if responseData.Category == "" {
			errors = append(errors, "Please enter item category")
		}
		if responseData.Price < 0 {
			errors = append(errors, "Please enter non-negative item price/value")
		}
		if len(errors) > 0 {
			templates["item_form"].Execute(writer, formDataForUserForItem{
				Item:   &responseData,
				Errors: errors,
			})
		} else {
			items = append(items, &responseData)
			templates["item_list"].Execute(writer, items)
		}

	}
}

func user_formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["user_form"].Execute(writer, formDataForUser{
			User: &User{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := User{
			Name:  request.Form["name"][0],
			Email: request.Form["email"][0],
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please enter user name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter email")
		}
		if len(errors) > 0 {
			templates["user_form"].Execute(writer, formDataForUser{
				User: &responseData, Errors: errors,
			})
		} else {
			users = append(users, &responseData)
			templates["user_list"].Execute(writer, users)
		}
	}
}

func addValues(x, y int) int {
	return x + y
}

// aboutHandler is the about page handler
func aboutHandler(writer http.ResponseWriter, request *http.Request) {
	sum := addValues(2, 3)

	n, err := fmt.Fprintf(writer, "<h1>About</h1>")
	_, _ = fmt.Fprintf(writer, fmt.Sprintf("The sum is %d", sum))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Number of bytes written: %d", n)
}

func main() {
	Populate()
	loadTemplates()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/user_form", user_formHandler)
	http.HandleFunc("/items", itemHandler)
	http.HandleFunc("/item_form", item_formHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Printf("Starting application on port %s", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
