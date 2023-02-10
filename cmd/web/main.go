package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/handlers"
	"github.com/gen4ralz/booking-app/internal/helpers"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/gen4ralz/booking-app/internal/render"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger  

func main(){
	
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr: port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
// what am I goinf to put in the session
gob.Register(models.Reservation{})
//check this to true when in production
app.InProduction = false

infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
app.InfoLog = infoLog

errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
app.ErrorLog = errorLog

session = scs.New()
session.Lifetime = 24 * time.Hour
session.Cookie.Persist = true // cookie persist until the browser window is closed by user.
session.Cookie.SameSite = http.SameSiteLaxMode
session.Cookie.Secure = app.InProduction // encrypt connection from Https, we don't need this in development

app.Session = session

tc, err := render.CreateTemplateCache()
if err != nil {
	log.Fatal("cannot create template cache")
	return err
}

app.TemplateCache = tc
app.UseCache = false

repo := handlers.NewRepo(&app)
handlers.NewHandler(repo)
//give render package can access to config
render.NewTemplates(&app)
helpers.NewHelpers(&app)

	return nil
}