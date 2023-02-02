package handlers

import (
	"net/http"

	"github.com/gen4ralz/booking-app/pkg/config"
	"github.com/gen4ralz/booking-app/pkg/models"
	"github.com/gen4ralz/booking-app/pkg/render"
)

// Repository  is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

// NewHandlers sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) HomePage(res http.ResponseWriter,req *http.Request){
	remoteIP := req.RemoteAddr
	m.App.Session.Put(req.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(res, "home.gohtml", &models.TemplateData{})
}

func (m *Repository) About(res http.ResponseWriter,req *http.Request){
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	// send the data to the template
	render.RenderTemplate(res, "about.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}