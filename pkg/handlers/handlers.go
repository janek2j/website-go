package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/janek2j/website-go/pkg/config"
	"github.com/janek2j/website-go/pkg/models"
	"github.com/janek2j/website-go/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig // holds the application config
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

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

type homeData struct {
	PresentTime   string
	PresentRandom int
}

type formDataForUser struct {
	*User
	Errors []string
}

type formDataForItem struct {
	*Item
	Errors []string
}

type AppUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var home = homeData{time.Now().String(), rand.Intn(6) + 1}

var users = make([]*User, 0, 10)
var items = make([]*Item, 0, 10)

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

func (m *Repository) HomeHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "home.page.html", &models.TemplateData{})
}

// AboutHandler is the about page handler
func (m *Repository) AboutHandler(writer http.ResponseWriter, request *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Testing template data"
	render.RenderTemplate(writer, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func UserHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "user_list.page.html", &models.TemplateData{})
}

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "item_list.page.html", &models.TemplateData{})
}

// func Item_formHandler(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method == http.MethodGet {
// 		render.Templates["item_form.page"].Execute(writer, formDataForItem{
// 			Item: &Item{}, Errors: []string{},
// 		})
// 	} else if request.Method == http.MethodPost {
// 		request.ParseForm()

// 		parsed_price, _ := strconv.ParseFloat(request.Form["price"][0], 32)
// 		responseData := Item{
// 			Name:        request.Form["name"][0],
// 			Category:    request.Form["category"][0],
// 			Description: request.Form["description"][0],
// 			Price:       float32(parsed_price),
// 			IsForRent:   request.Form["isForRent"][0] == "true",
// 		}

// 		errors := []string{}
// 		if responseData.Name == "" {
// 			errors = append(errors, "Please enter item name")
// 		}
// 		if responseData.Description == "" {
// 			errors = append(errors, "Please enter some description")
// 		}
// 		if responseData.Category == "" {
// 			errors = append(errors, "Please enter item category")
// 		}
// 		if responseData.Price < 0 {
// 			errors = append(errors, "Please enter non-negative item price/value")
// 		}
// 		if len(errors) > 0 {
// 			render.Templates["item_form.page"].Execute(writer, formDataForItem{
// 				Item:   &responseData,
// 				Errors: errors,
// 			})
// 		} else {
// 			items = append(items, &responseData)
// 			render.Templates["item_list.page"].Execute(writer, items)
// 		}

// 	}
// }

// func User_formHandler(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method == http.MethodGet {
// 		render.Templates["user_form.page"].Execute(writer, formDataForUser{
// 			User: &User{}, Errors: []string{},
// 		})
// 	} else if request.Method == http.MethodPost {
// 		request.ParseForm()
// 		responseData := User{
// 			Name:  request.Form["name"][0],
// 			Email: request.Form["email"][0],
// 		}

// 		errors := []string{}
// 		if responseData.Name == "" {
// 			errors = append(errors, "Please enter user name")
// 		}
// 		if responseData.Email == "" {
// 			errors = append(errors, "Please enter email")
// 		}
// 		if len(errors) > 0 {
// 			render.Templates["user_form"].Execute(writer, formDataForUser{
// 				User: &responseData, Errors: errors,
// 			})
// 		} else {
// 			users = append(users, &responseData)
// 			render.Templates["user_list"].Execute(writer, users)
// 		}
// 	}
// }

func addValues(x, y int) int {
	return x + y
}

func OtherHandler(writer http.ResponseWriter, request *http.Request) {
	sum := addValues(2, 3)

	n, err := fmt.Fprintf(writer, "<h1>About</h1>")
	_, _ = fmt.Fprintf(writer, fmt.Sprintf("The sum is %d", sum))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Number of bytes written: %d", n)
}
