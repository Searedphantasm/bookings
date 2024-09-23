package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Searedphantasm/bookings/internal/config"
	"github.com/Searedphantasm/bookings/internal/driver"
	"github.com/Searedphantasm/bookings/internal/forms"
	"github.com/Searedphantasm/bookings/internal/helpers"
	"github.com/Searedphantasm/bookings/internal/models"
	"github.com/Searedphantasm/bookings/internal/render"
	"github.com/Searedphantasm/bookings/internal/repository"
	"github.com/Searedphantasm/bookings/internal/repository/dbrepo"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (m *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "about.page.gohtml", &models.TemplateData{})
}

// Reservations renders the make a reservations and display form
func (m *Repository) Reservations(writer http.ResponseWriter, request *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(writer, request, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservations handles the posting of reservation form.
func (m *Repository) PostReservations(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation := models.Reservation{
		FirstName: request.Form.Get("first_name"),
		LastName:  request.Form.Get("last_name"),
		Phone:     request.Form.Get("phone"),
		Email:     request.Form.Get("email"),
	}

	form := forms.New(request.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(writer, request, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(request.Context(), "reservation", reservation)

	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "generals.page.gohtml", &models.TemplateData{})
}

// Majors is the about majors room handler
func (m *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the availability page handler
func (m *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability
func (m *Repository) PostAvailability(writer http.ResponseWriter, request *http.Request) {
	start := request.Form.Get("start")
	end := request.Form.Get("end")

	writer.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson handles request for availability and send JSON response
func (m *Repository) AvailabilityJson(writer http.ResponseWriter, request *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	//log.Println(string(out))
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(out)
}

// Contact is the contact page handler
func (m *Repository) Contact(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := m.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(request.Context(), "error", "Can't get reservation from session")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(request.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(writer, request, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}
