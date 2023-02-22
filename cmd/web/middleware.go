package main

import (
	"net/http"

	"github.com/gen4ralz/booking-app/internal/helpers"
	"github.com/justinas/nosurf"
)

// Use for test middleware
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(res, req)
// 	})
// }

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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !helpers.IsAuthenticate(req){
			session.Put(req.Context(), "error", "Log in first")
			http.Redirect(res,req,"/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(res,req)
	})
}