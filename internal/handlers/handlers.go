package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chelobotix/booking-go/internal/config"
	"github.com/chelobotix/booking-go/internal/forms"
	"github.com/chelobotix/booking-go/internal/helpers"
	"github.com/chelobotix/booking-go/internal/models"
	"github.com/chelobotix/booking-go/internal/render"
	"net/http"
)

// Repo the repository using by the handlers
var Repo *Repository

// Repository ins the repository type
type Repository struct {
	AppConfig *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: appConfig,
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(repository *Repository) {
	Repo = repository
}

// Home is the handler for the home page
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Generals is the handler for the home page
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Reservations is the handler for the home page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservations is the handler for the home page
func (repo *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	err = errors.New("This is an error")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	repo.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.AppConfig.ErrorLog.Println("Cannot get error from session")
		repo.AppConfig.Session.Put(r.Context(), "error", "cannon get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.AppConfig.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

// Major is the handler for the home page
func (repo *Repository) Major(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the handler for the home page
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability is the handler for the home page
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON is the handler for the home page
func (repo *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	response := jsonResponse{
		Ok:      true,
		Message: "Available",
	}

	outResponse, err := json.MarshalIndent(response, "", "")
	if err != nil {
		helpers.ServerError(w, err)
	}
	w.Header().Set("Content-Type", "application-json")
	w.Write(outResponse)
}

// Contact is the handler for the home page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}
