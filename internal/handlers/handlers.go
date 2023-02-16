package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/driver"
	"github.com/gen4ralz/booking-app/internal/forms"
	"github.com/gen4ralz/booking-app/internal/helpers"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/gen4ralz/booking-app/internal/render"
	"github.com/gen4ralz/booking-app/internal/repository"
	"github.com/gen4ralz/booking-app/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App:a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandler(r *Repository){
	Repo = r
}

func (m *Repository)HomePage(res http.ResponseWriter, req *http.Request){ 
	render.Template(res,req, "home.gohtml", &models.TemplateData{})
}

func (m *Repository)About(res http.ResponseWriter, req *http.Request) {

	render.Template(res,req, "about.gohtml", &models.TemplateData{})
}

func (m *Repository) Generals (res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "generals.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors (res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "majors.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation (res http.ResponseWriter,req *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(res,req, "reservation.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation (res http.ResponseWriter,req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: req.FormValue("first_name"),
		LastName:  req.FormValue("last_name"),
		Email:     req.FormValue("email"),
		Phone:     req.FormValue("phone"),
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, req)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(res, req, "reservation.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(req.Context(), "reservation", reservation)
	http.Redirect(res, req, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Availability (res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "search-availability.gohtml", &models.TemplateData{})
}

type jsonResponse struct {
	OK 			bool		`json:"ok"`
	Message 	string		`json:"message"`
}

func (m *Repository) AvailabilityJSON (res http.ResponseWriter,req *http.Request) {
	resp := jsonResponse {
		OK: true,
		Message: "Available!",
	}
	json,err := json.MarshalIndent(resp,"","     ")
	if err != nil {
		helpers.ServerError(res, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
}

func (m *Repository) PostAvailability (res http.ResponseWriter,req *http.Request) {
	start := req.FormValue("start")
	end := req.FormValue("end")
	res.Write([]byte(fmt.Sprintf("start date is %s, end date is %s", start, end)))
}

func (m *Repository) Contact (res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "contact.gohtml", &models.TemplateData{})
}

func (m *Repository) ReservationSummary (res http.ResponseWriter,req *http.Request) {
	// models.Reservation is a type
	reservation,ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(req.Context(), "error", "Can't get reservation from session")
		http.Redirect(res,req,"/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(req.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(res,req, "reservation-summary.gohtml", &models.TemplateData{
		Data: data,
	})
}