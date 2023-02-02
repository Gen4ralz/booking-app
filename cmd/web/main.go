package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gen4ralz/booking-app/pkg/config"
	"github.com/gen4ralz/booking-app/pkg/handlers"
	"github.com/gen4ralz/booking-app/pkg/render"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager  

func main(){

	//check this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // cookie persist until the browser window is closed by user.
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // encrypt connection from Https, we don't need this in development

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	//give render package can access to config
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr: port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}