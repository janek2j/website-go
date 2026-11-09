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

// CfgRepo is the repository used by the handlers
var CfgRepo *ConfigRepository

// ConfigRepository is the repository type
type ConfigRepository struct {
	App *config.AppConfig // holds the application config
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *ConfigRepository {
	return &ConfigRepository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *ConfigRepository) {
	CfgRepo = r
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

func (m *ConfigRepository) HomeHandler(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	// req := strconv.FormatInt(r.ContentLength+10, 10)
	req := string(r.ContentLength + 74)
	fmt.Println("req:", req)

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	m.App.Session.Put(r.Context(), "req", req)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{
		StringMap: map[string]string{
			"PresentRandom": "999999",
		},
	})
}

// AboutHandler is the about page handler
func (m *ConfigRepository) AboutHandler(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Testing template data"

	stringMap["remote_ip"] = m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["req"] = m.App.Session.GetString(r.Context(), "req")

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "user_list.page.html", &models.TemplateData{})
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "item_list.page.html", &models.TemplateData{})
}

func RentalHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "rental_list.page.html", &models.TemplateData{})
}

// func Item_formHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		render.Templates["item_form.page"].Execute(w, formDataForItem{
// 			Item: &Item{}, Errors: []string{},
// 		})
// 	} else if r.Method == http.MethodPost {
// 		r.ParseForm()

// 		parsed_price, _ := strconv.ParseFloat(r.Form["price"][0], 32)
// 		responseData := Item{
// 			Name:        r.Form["name"][0],
// 			Category:    r.Form["category"][0],
// 			Description: r.Form["description"][0],
// 			Price:       float32(parsed_price),
// 			IsForRent:   r.Form["isForRent"][0] == "true",
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
// 			render.Templates["item_form.page"].Execute(w, formDataForItem{
// 				Item:   &responseData,
// 				Errors: errors,
// 			})
// 		} else {
// 			items = append(items, &responseData)
// 			render.Templates["item_list.page"].Execute(w, items)
// 		}

// 	}
// }

// func User_formHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		render.Templates["user_form.page"].Execute(w, formDataForUser{
// 			User: &User{}, Errors: []string{},
// 		})
// 	} else if r.Method == http.MethodPost {
// 		r.ParseForm()
// 		responseData := User{
// 			Name:  r.Form["name"][0],
// 			Email: r.Form["email"][0],
// 		}

// 		errors := []string{}
// 		if responseData.Name == "" {
// 			errors = append(errors, "Please enter user name")
// 		}
// 		if responseData.Email == "" {
// 			errors = append(errors, "Please enter email")
// 		}
// 		if len(errors) > 0 {
// 			render.Templates["user_form"].Execute(w, formDataForUser{
// 				User: &responseData, Errors: errors,
// 			})
// 		} else {
// 			users = append(users, &responseData)
// 			render.Templates["user_list"].Execute(w, users)
// 		}
// 	}
// }

func addValues(x, y int) int {
	return x + y
}

func OtherHandler(w http.ResponseWriter, r *http.Request) {
	sum := addValues(2, 3)

	n, err := fmt.Fprintf(w, "<h1>About</h1>")
	_, _ = fmt.Fprintf(w, fmt.Sprintf("The sum is %d", sum))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Number of bytes written: %d", n)
}
