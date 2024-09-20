package handlers

import (
	"github.com/Searedphantasm/bookings/pkg/config"
	"github.com/Searedphantasm/bookings/pkg/models"
	"github.com/Searedphantasm/bookings/pkg/render"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (m *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "about.page.gohtml", &models.TemplateData{})
}

// Reservations renders the make a reservations and display form
func (m *Repository) Reservations(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Generals renders the room page
func (m *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "generals.page.gohtml", &models.TemplateData{})
}

// Majors is the about majors room handler
func (m *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the availability page handler
func (m *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "search-availability.page.gohtml", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "contact.page.gohtml", &models.TemplateData{})
}
