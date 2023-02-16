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
	"github.com/gen4ralz/booking-app/internal/driver"
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
	
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr: port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB,error) {
// what am I goinf to put in the session
gob.Register(models.Reservation{})
gob.Register(models.User{})
gob.Register(models.Room{})
gob.Register(models.Restriction{})
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

//connect to database
log.Println("Connecting to database...")
db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=nattapongpanyaakratham password=")
if err != nil {
	log.Fatal("Cannot connect to database!")
}
log.Println("Connected to database!")

tc, err := render.CreateTemplateCache()
if err != nil {
	log.Fatal("cannot create template cache")
	return nil, err
}

app.TemplateCache = tc
app.UseCache = false

repo := handlers.NewRepo(&app, db)
handlers.NewHandler(repo)
//give render package can access to config
render.NewRenderer(&app)
helpers.NewHelpers(&app)

	return db, nil
}