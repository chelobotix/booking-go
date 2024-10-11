package render

import (
	"bytes"
	"github.com/chelobotix/booking-go/pkg/config"
	"github.com/chelobotix/booking-go/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates set the config fot the template package
func NewTemplates(appConfig *config.AppConfig) {
	app = appConfig
}

var templateCache map[string]*template.Template

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	if app.UseCache {
		// Create a template cache
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := templateCache[tmpl]

	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buffer := new(bytes.Buffer)
	td = AddDefaultData(td)

	_ = t.Execute(buffer, td)

	// render the template
	_, err := buffer.WriteTo(w)

	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.gohtml")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.gohtml")

		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, err
}
