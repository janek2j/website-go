package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	fmt.Println("Template data:\n", td)
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, templateName string, tmplData *models.TemplateData) {
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
	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, r)

	_ = tmpl.Execute(buf, tmplData)

	// render the template in a browser
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser:", err)
	}
}

// CreateTemplateCache creates template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
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

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return tmplCache, err
		}

		// combines .page template with .layout template together
		if len(matches) > 0 {
			currentTemplate, err = currentTemplate.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				fmt.Println(err)
				return tmplCache, err
			}
		}

		tmplCache[name] = currentTemplate
	}
	return tmplCache, nil
}
