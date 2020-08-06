package main

import (
	"context"
	"fmt"
	"net/http"

	"karolharasim.com/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		//Check if a userID value exists in the session. if it isn't, call the other handler in chain as normal
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		//Fetch the details of current user from the database. If no matching record is found, remove the userID from their session and call the next handler
		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return 
		}

		//Otherwise we know, that request is coming from a valid authenticated, logged in user. Create a new copy of the request with the user info added to it
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}



func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s-%s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//deferred func that will always run in case of panic
		defer func() {
			//use the builtin recover function ythat checks, if there has been a panic or not
			if err := recover(); err != nil {
				//set a connection-close header
				w.Header().Set("Connection", "Close")
				//call the app.ServerError method to return a 500 internal error response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If the user is not authenticated, redirect them to the login page and return from the middleware chain so that no subsequent handlers in the chain are executed
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", 302)
			return
		}

		//Otherwise call the next handler in the chain
		next.ServeHTTP(w, r)

	})
}

//middleware function which uses a customized CSRF cookie with the Secure, Path and HTTPOnly flags set
func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie {
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})

	return csrfHandler
}