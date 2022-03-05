package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/erdauletbatalov/forum.git/pkg/models"
	"github.com/erdauletbatalov/forum.git/pkg/session"
	"github.com/google/uuid"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		isSession, user := session.IsSession(r)

		app.render(w, r, "home.page.html", &templateData{
			IsSession: isSession, //
			User:      user,
		})
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("POST!!!")
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Username: r.FormValue("nickname"),
		}
		err := app.forum.CreateUser(&user)
		if err != nil {
			switch err.Error() {
			case "UNIQUE constraint failed: users.email":
				app.render(w, r, "signup.page.html", &templateData{
					IsError: isError{true, "this email is already in use"},
				})
				return
			case "UNIQUE constraint failed: users.nickname":
				app.render(w, r, "signup.page.html", &templateData{
					IsError: isError{true, "this nickname is already in use"},
				})
				return

			}
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
	case http.MethodGet:
		app.render(w, r, "signup.page.html", &templateData{})
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		// foundUser := &models.User{}
		var err error
		err = app.forum.LogInUser(&user)
		if err != nil {
			app.render(w, r, "signin.page.html", &templateData{
				IsError: isError{true, "incorrect email or password"},
			})
			return
		}
		// Create a new random session token
		// we use the "github.com/google/uuid" library to generate UUIDs
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		// Set the token in the session map, along with the session information
		session.Sessions[sessionToken] = session.Session{
			Email:  user.Email,
			Expiry: expiresAt,
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case http.MethodGet:
		app.render(w, r, "signin.page.html", &templateData{})
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/profile" {
		app.notFound(w)
		return
	}
	switch r.Method {
	case http.MethodGet:
		isSession, user := session.IsSession(r)
		if isSession {
			var err error
			user, err = app.forum.GetUser(user.Email)
			if err != nil {
				fmt.Println("error email not found")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			app.render(w, r, "profile.page.html", &templateData{
				IsSession: isSession,
				User:      user,
			})
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		isSession, _ := session.IsSession(r)
		if isSession {
			c, _ := r.Cookie("session_token")
			sessionToken := c.Value
			delete(session.Sessions, sessionToken)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}
