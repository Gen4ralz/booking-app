package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gen4ralz/booking-app/internal/config"
)

var app *config.AppConfig

// NewHelpers set up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

//use responseWriter because we need something to write to client
func ClientError(res http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(res, http.StatusText(status), status)
}

func ServerError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s",err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}