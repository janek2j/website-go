package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const portNumber = ":5000" // in ":nnnn" format

type Rsvp struct {
	Name, Email string
	IsActive    bool
}

type Deck struct {
	Name        string
	Description string
	IsActive    bool
}

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
	Level uint8  `json:"level"`
	Deck  string `json:"deck"`
}

type formData struct {
	*Deck
	Errors []string
}

type welcomeData struct {
	PresentTime   string
	PresentRandom int
}

type AppUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var decks = make([]*Deck, 0, 10)
var cards = make([]*Card, 0, 10)

var welcome = welcomeData{time.Now().String(), rand.Intn(6) + 1}

var templates = make(map[string]*template.Template, 3)

// Populate fills Deck and Card collections with some initial data
func Populate() {
	decks = append(decks,
		&Deck{"Esperanto", "słówka", true},
		&Deck{"Dates nnnn", "daty historyczne", false},
	)

	jsonForCards := `
	[
		{
			"front": "kot",
			"back": "kato",
			"level": 0,
			"deck": "1"
		},
		{
			"front": "pies",
			"back": "hundo",
			"level": 0,
			"deck": "1"
		},
		{
			"front": "chleb",
			"back": "pano",
			"level": 0,
			"deck": "1"
		}
	]`

	var unmarshalledCards []Card
	err := json.Unmarshal([]byte(jsonForCards), &unmarshalledCards)
	if err != nil {
		log.Println("Error unmarshalling card json", err)
	}
	log.Printf("unmarshalled: %v", unmarshalledCards)
	for _, c := range unmarshalledCards {
		cards = append(cards, &c)
	}

}

func loadTemplates() {
	templateNames := [4]string{"welcome", "deck_list", "card_list", "deck_form"}
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

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, welcome)
}

func deckHandler(writer http.ResponseWriter, request *http.Request) {
	templates["deck_list"].Execute(writer, decks)
}

func cardHandler(writer http.ResponseWriter, request *http.Request) {
	templates["card_list"].Execute(writer, cards)
}

func deck_formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["deck_form"].Execute(writer, formData{
			Deck: &Deck{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Deck{
			Name:        request.Form["name"][0],
			Description: request.Form["description"][0],
			IsActive:    request.Form["isactive"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please enter deck name")
		}
		if responseData.Description == "" {
			errors = append(errors, "Please enter deck description")
		}
		if len(errors) > 0 {
			templates["deck_form"].Execute(writer, formData{
				Deck: &responseData, Errors: errors,
			})
		} else {
			decks = append(decks, &responseData)
			templates["deck_list"].Execute(writer, decks)
			// if responseData.IsActive {
			// 	templates["thanks"].Execute(writer, responseData.Name)
			// } else {
			// 	templates["sorry"].Execute(writer, responseData.Name)
			// }
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

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/decks", deckHandler)
	http.HandleFunc("/cards", cardHandler)
	http.HandleFunc("/deck_form", deck_formHandler)
	http.HandleFunc("/about", aboutHandler)

	// http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
	// 	n, err := fmt.Fprintf(w, "<h1>About</h1>")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	log.Printf("Number of bytes written: %d", n)
	// })

	fmt.Printf("Starting application on port %s", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
