package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig


func TestMain(m *testing.M) {

	// what am I goinf to put in the session
gob.Register(models.Reservation{})
//check this to true when in production
testApp.InProduction = false

infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
testApp.InfoLog = infoLog

errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
testApp.ErrorLog = errorLog

session = scs.New()
session.Lifetime = 24 * time.Hour
session.Cookie.Persist = true // cookie persist until the browser window is closed by user.
session.Cookie.SameSite = http.SameSiteLaxMode
session.Cookie.Secure = false // encrypt connection from Https, we don't need this in development

testApp.Session = session

app = &testApp

	os.Exit(m.Run())
}

type myWriter struct {}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte)(int, error){
	length := len(b)
	return length, nil
}