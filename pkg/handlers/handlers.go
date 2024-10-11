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

	render.RenderTemplate(w, "home.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello buddy!"

	remoteIP := repo.AppConfig.Session.GetString(r.Context(), "remote-ip")
	stringMap["remoteIP`"] = remoteIP

	render.RenderTemplate(w, "about.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
