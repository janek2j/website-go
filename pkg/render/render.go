package render

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/janek2j/website-go/pkg/config"
	"github.com/janek2j/website-go/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, templateName string, tmplData *models.TemplateData) {
	var tmplCache map[string]*template.Template
	var err error
	if app.UseCache {
		// get the template cache from the app config
		tmplCache = app.TemplateCache
	} else {
		// create a new template cache
		tmplCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("cannot create template cache:", err)
		}
	}

	tmpl, ok := tmplCache[templateName]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	// for fine grained error checking
	// buf := new(bytes.Buffer)
	// err = tmpl.Execute(buf, nil)
	// if err != nil {
	// 	log.Println(err)
	// }

	// render the template
	// _, err = buf.WriteTo(w)
	// if err != nil {
	// fmt.Println("error writing template to browser:", err)
	// }

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		fmt.Println(err)
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/"+templateName+".html", "./templates/base.layout.html")
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	log.Println("error passing template:", err)
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return tmplCache, err
	}

	// range through all files ending with *page.html
	for _, page := range pages {
		name := filepath.Base(page)
		currentTemplate, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return tmplCache, err
		}

		// combines .page template with .layout template together
		if len(matches) > 0 {
			currentTemplate, err = currentTemplate.ParseGlob("./templates/*.layout.html")
			if err != nil {
				fmt.Println(err)
				return tmplCache, err
			}
		}

		tmplCache[name] = currentTemplate
	}
	return tmplCache, nil
}
