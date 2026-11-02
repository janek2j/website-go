package render

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var Templates = make(map[string]*template.Template, 3)

func LoadTemplates() {
	templateNames := [6]string{"home", "user_list", "item_list", "user_form", "item_form", "about"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("./templates/layout.html", "./templates/"+name+".html")
		if err == nil {
			Templates[name] = t
			log.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

var templateCache = make(map[string]*template.Template)

func RenderTemplateSecond(w http.ResponseWriter, templateName string) {
	var tmpl *template.Template
	var err error

	// check to see if a template is already in the cache
	_, inMap := templateCache[templateName]
	if !inMap {
		// need to create the template
		log.Println("creating template and adding to cache")
		err = createTemplateCache(templateName)
		if err != nil {
			log.Println(err)
		}
	} else {
		// the template is already in the cache
		log.Println("using cached template")
	}

	tmpl = templateCache[templateName]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(templateName string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", templateName),
		"./templates/layout.html",
	}

	// parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache (map)
	templateCache[templateName] = tmpl

	return nil
}
