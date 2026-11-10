package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/driver"
	"github.com/janek2j/website-go/internal/forms"
	"github.com/janek2j/website-go/internal/models"
	"github.com/janek2j/website-go/internal/render"
	"github.com/janek2j/website-go/internal/repository"
	"github.com/janek2j/website-go/internal/repository/dbrepo"
)

// CfgRepo is the repository used by the handlers
var CfgRepo *ConfigRepository

// ConfigRepository is the repository type
type ConfigRepository struct {
	App *config.AppConfig // holds the application config
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *ConfigRepository {
	return &ConfigRepository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
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

// type formDataForUser struct {
// 	*User
// 	Errors []string
// }

// type formDataForItem struct {
// 	*Item
// 	Errors []string
// }

type AppUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

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

func (m *ConfigRepository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	// req := strconv.FormatInt(r.ContentLength+10, 10)
	req := string(r.ContentLength + 74)
	fmt.Println("req:", req)

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	m.App.Session.Put(r.Context(), "req", req)

	render.Template(w, r, "home.page.html", &models.TemplateData{
		StringMap: map[string]string{},
	})
}

// AboutHandler is the about page handler
func (m *ConfigRepository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Testing template data"

	stringMap["remote_ip"] = m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["req"] = m.App.Session.GetString(r.Context(), "req")

	render.Template(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *ConfigRepository) Book(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "book.page.html", &models.TemplateData{})
}

func (m *ConfigRepository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostSearchAvailability renders the search availability page
func (m *ConfigRepository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	fmt.Fprintf(w, "start date is %s and end date is %s", start, end)
}

type responseJSON struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *ConfigRepository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := responseJSON{
		OK:      true,
		Message: "Dostępne!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		fmt.Println("error when marshalling to json:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// MakeReservation renders the make a reservation page and displays form
func (m *ConfigRepository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostMakeReservation handles the posting of reservation form
func (m *ConfigRepository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmailValidator("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

	// w.Write([]byte("Posted to make-reservation"))
	// render.Template(w, "", nil)
}

// Contact renders the page with contact info
func (m *ConfigRepository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *ConfigRepository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	// data := m.App.Session.Pop(r.Context(), "reservation")
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	m.App.Session.Remove(r.Context(), "reservation")

	if !ok {
		log.Println("cannot get info from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}

func UserH(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "user_list.page.html", &models.TemplateData{})
}

func ItemH(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "item_list.page.html", &models.TemplateData{})
}

func Rental(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "rental_list.page.html", &models.TemplateData{})
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
