package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/helpers"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/gen4ralz/booking-app/internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app config.AppConfig
var session *scs.SessionManager
// var infoLog *log.Logger
// var errorLog *log.Logger 
var pathToTemplate = "./../../templates"

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})
	//check this to true when in production
	app.InProduction = false
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // cookie persist until the browser window is closed by user.
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // encrypt connection from Https, we don't need this in development
	
	app.Session = session
	
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	
	app.TemplateCache = tc
	app.UseCache = true
	
	repo := NewRepo(&app)
	NewHandler(repo)
	//give render package can access to config
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	mux := chi.NewRouter();

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.HomePage)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	
	mux.Get("/contact", Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction, // In development mode we don't run in Https rightnow.
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache()(map[string]*template.Template, error){
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.gohtml from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", pathToTemplate))
	if err != nil {
		return myCache,err
	}

	// range through all files ending with *.gohtml
	for _,page := range pages {
		name := filepath.Base(page)

		//ts = template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}