package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/driver"
	"github.com/gen4ralz/booking-app/internal/forms"
	"github.com/gen4ralz/booking-app/internal/helpers"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/gen4ralz/booking-app/internal/render"
	"github.com/gen4ralz/booking-app/internal/repository"
	"github.com/gen4ralz/booking-app/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
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
	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(res,errors.New("cannot get reservation from session"))
		return
	}

	room,err := m.DB.GetRoomById(reservation.RoomID)
	if err != nil {
		helpers.ServerError(res,err)
		return
	}

	reservation.Room.RoomName = room.RoomName

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap  := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(res,req, "reservation.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation (res http.ResponseWriter,req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	sd := req.FormValue("start_date")
	ed := req.FormValue("end_date")

	// 01-01-2020
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	roomID, err := strconv.Atoi(req.FormValue("room_id"))
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	reservation := models.Reservation{
		FirstName: req.FormValue("first_name"),
		LastName:  req.FormValue("last_name"),
		Email:     req.FormValue("email"),
		Phone:     req.FormValue("phone"),
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomID,
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
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

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(res, err)
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

		// 01-01-2020
		layout := "2006-01-02"

		startDate, err := time.Parse(layout, start)
		if err != nil {
			helpers.ServerError(res, err)
			return
		}
	
		endDate, err := time.Parse(layout, end)
		if err != nil {
			helpers.ServerError(res, err)
			return
		}

		rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
		if err != nil {
			helpers.ServerError(res, err)
			return
		}

		for _, i := range rooms {
			m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
		}

		if len(rooms) == 0 {
			// no available
			m.App.Session.Put(req.Context(), "error", "No availability")
			http.Redirect(res, req, "/search-availability", http.StatusSeeOther)
			return
		}

		data := make(map[string]interface{})
		data["rooms"] = rooms

		reservation := models.Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}

		m.App.Session.Put(req.Context(),"reservation", reservation)

		render.Template(res,req, "choose-room.gohtml", &models.TemplateData{
			Data: data,
		})
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

func (m *Repository) ChooseRoom(res http.ResponseWriter,req *http.Request){
	roomID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(res, err)
		return
	}

	reservation.RoomID = roomID

	m.App.Session.Put(req.Context(), "reservation", reservation)

	http.Redirect(res, req, "/make-reservation", http.StatusSeeOther)
}