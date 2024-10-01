package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Searedphantasm/bookings/internal/config"
	"github.com/Searedphantasm/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}
var app *config.AppConfig
var pathToTemplates = "./templates"

func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// NewRenderer sets the config for the template pkg
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate returns time with specified format
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultData add default data to all templates
func AddDefaultData(td *models.TemplateData, request *http.Request) *models.TemplateData {
	// put something in the session , until the next time a page is displayed, and then it's taken out.
	td.Flash = app.Session.PopString(request.Context(), "flash")
	td.Error = app.Session.PopString(request.Context(), "error")
	td.Warning = app.Session.PopString(request.Context(), "warning")
	td.CSRFToken = nosurf.Token(request)
	if app.Session.Exists(request.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template renders html templates
func Template(writer http.ResponseWriter, request *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
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
		return err
	}

	return nil
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
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
