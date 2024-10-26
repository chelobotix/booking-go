package handlers

import (
	"encoding/json"
	"errors"
	"github.com/chelobotix/booking-go/internal/config"
	"github.com/chelobotix/booking-go/internal/driver"
	"github.com/chelobotix/booking-go/internal/forms"
	"github.com/chelobotix/booking-go/internal/helpers"
	"github.com/chelobotix/booking-go/internal/models"
	"github.com/chelobotix/booking-go/internal/render"
	"github.com/chelobotix/booking-go/internal/repository"
	"github.com/chelobotix/booking-go/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository using by the handlers
var Repo *Repository

// Repository ins the repository type
type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConfig: appConfig,
		DB:        dbrepo.NewPostgresRepo(db.SQL, appConfig),
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(repository *Repository) {
	Repo = repository
}

// Home is the handler for the home page
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Generals is the handler for the home page
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Reservations is the handler for the home page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err := repo.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	res.Room.RoomName = room.RoomName

	startDate := res.StartDate.Format("2006-01-02")
	endDate := res.StartDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservations is the handler for the home page
func (repo *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationId, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
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
	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

// Major is the handler for the home page
func (repo *Repository) Major(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the handler for the home page
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability is the handler for the home page
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endtDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
	}

	availableRooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endtDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(availableRooms) == 0 {
		repo.AppConfig.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["availableRooms"] = availableRooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endtDate,
	}

	repo.AppConfig.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
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

	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := repo.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomId

	repo.AppConfig.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
