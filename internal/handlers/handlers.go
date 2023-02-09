package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/gen4ralz/booking-app/internal/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{a}
}

func NewHandler(r *Repository){
	Repo = r
}

func (m *Repository)HomePage(res http.ResponseWriter, req *http.Request){ 
	remoteIp := req.RemoteAddr
	m.App.Session.Put(req.Context(), "remote_ip", remoteIp)
	render.RenderTemplate(res,req, "home.gohtml", &models.TemplateData{})
}

func (m *Repository)About(res http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp


	render.RenderTemplate(res,req, "about.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Generals (res http.ResponseWriter,req *http.Request) {
	render.RenderTemplate(res,req, "generals.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors (res http.ResponseWriter,req *http.Request) {
	render.RenderTemplate(res,req, "majors.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation (res http.ResponseWriter,req *http.Request) {
	render.RenderTemplate(res,req, "reservation.gohtml", &models.TemplateData{})
}

func (m *Repository) Availability (res http.ResponseWriter,req *http.Request) {
	render.RenderTemplate(res,req, "search-availability.gohtml", &models.TemplateData{})
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
		log.Println(err)
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
	render.RenderTemplate(res,req, "contact.gohtml", &models.TemplateData{})
}