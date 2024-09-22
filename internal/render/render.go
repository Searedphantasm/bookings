package render

import (
	"bytes"
	"fmt"
	"github.com/Searedphantasm/bookings/internal/config"
	"github.com/Searedphantasm/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates sets the cofig for the template pkg
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData add default data to all templates
func AddDefaultData(td *models.TemplateData, request *http.Request) *models.TemplateData {
	// put something in the session , until the next time a page is displayed, and then it's taken out.
	td.Flash = app.Session.PopString(request.Context(), "flash")
	td.Error = app.Session.PopString(request.Context(), "error")
	td.Warning = app.Session.PopString(request.Context(), "warning")
	td.CSRFToken = nosurf.Token(request)
	return td
}

// RenderTemplate renders html templates
func RenderTemplate(writer http.ResponseWriter, request *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error loading template:", tmpl)
	}

	buf := new(bytes.Buffer) // something that holds bytes

	td = AddDefaultData(td, request)

	err := t.Execute(buf, td)

	if err != nil {
		log.Println("Could not get template from template cache")
	}

	// render the template
	_, err = buf.WriteTo(writer)
	if err != nil {
		log.Println(err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	// range through all files ending with *.page.gohtml
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
