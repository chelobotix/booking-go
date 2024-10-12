package handlers

import (
	"github.com/chelobotix/booking-go/pkg/config"
	"github.com/chelobotix/booking-go/pkg/models"
	"github.com/chelobotix/booking-go/pkg/render"
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
	remoteIP := r.RemoteAddr
	repo.AppConfig.Session.Put(r.Context(), "remote-ip", remoteIP)

	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello buddy!"

	remoteIP := repo.AppConfig.Session.GetString(r.Context(), "remote-ip")
	stringMap["remoteIP`"] = remoteIP

	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the handler for the home page
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "generals.page.gohtml", &models.TemplateData{})
}

// Reservations is the handler for the home page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Major is the handler for the home page
func (repo *Repository) Major(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the handler for the home page
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability is the handler for the home page
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Posted to search availability"))
}

// Contact is the handler for the home page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "contact.page.gohtml", &models.TemplateData{})
}
