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

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App:a,
		DB: dbrepo.NewTestingsRepo(a),
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


// Reservation
func (m *Repository) Reservation (res http.ResponseWriter,req *http.Request) {
	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(req.Context(), "error", "can't get reservation from session")
		http.Redirect(res,req, "/", http.StatusTemporaryRedirect)
		return
	}

	room,err := m.DB.GetRoomById(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't find room")
		http.Redirect(res,req, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room.RoomName = room.RoomName

	m.App.Session.Put(req.Context(), "reservation", reservation)

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

//PostReservation
func (m *Repository) PostReservation (res http.ResponseWriter,req *http.Request) {
	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(res, errors.New("can't get from session"))
		return
	}

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation.FirstName = req.FormValue("first_name")
	reservation.LastName = req.FormValue("last_name")
	reservation.Phone = req.FormValue("phone")
	reservation.Email = req.FormValue("email")

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
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
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
	RoomID		string		`json:"room_id"`
	StartDate	string		`json:"start_date"`
	EndDate		string		`json:"end_date"`
}

func (m *Repository) AvailabilityJSON (res http.ResponseWriter,req *http.Request) {

	sd := req.Form.Get("start")
	ed := req.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID,_ := strconv.Atoi(req.Form.Get("room_id"))
	
	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse {
		OK: available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
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

// ReservationSummary
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(res,req, "reservation-summary.gohtml", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
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

func (m *Repository) BookRoom(res http.ResponseWriter, req *http.Request) {
	// id, s, e

	//convert from string to int
	roomID, _ := strconv.Atoi(req.URL.Query().Get("id"))
	sd := req.URL.Query().Get("s")
	ed := req.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var reservation models.Reservation

	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(res, err)
		return
	}

	reservation.Room.RoomName = room.RoomName

	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	m.App.Session.Put(req.Context(), "reservation", reservation)

	log.Println(roomID, startDate, endDate)
	http.Redirect(res,req, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Login (res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "login.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(res http.ResponseWriter,req *http.Request) {
	// Anytime you doing login or logout make sure you renew the token
	_ = m.App.Session.RenewToken(req.Context())

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	form := forms.New(req.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	
	if !form.Valid() {
		// take user back to page
		render.Template(res,req,"login.gohtml", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(req.Context(), "error", "Invalid login credentials")
		http.Redirect(res,req, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(req.Context(), "user_id", id)
	m.App.Session.Put(req.Context(), "flash", "Logged on successfully")
	http.Redirect(res,req, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(res http.ResponseWriter, req *http.Request){
	_ = m.App.Session.Destroy(req.Context())
	_ = m.App.Session.RenewToken(req.Context())

	http.Redirect(res,req,"/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(res http.ResponseWriter,req *http.Request) {
	render.Template(res,req, "admin-dashboard.gohtml", &models.TemplateData{})
}