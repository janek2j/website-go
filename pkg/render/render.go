package render

import (
	"log"
	"text/template"
)

var Templates = make(map[string]*template.Template, 3)

func LoadTemplates() {
	templateNames := [6]string{"home", "user_list", "item_list", "user_form", "item_form", "about"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("./templates/"+"layout.html", "./templates/"+name+".html")
		if err == nil {
			Templates[name] = t
			log.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}
