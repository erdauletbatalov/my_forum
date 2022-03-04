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
			IsSession: isSession,
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
			Nickname: r.FormValue("nickname"),
		}
		err := app.forum.CreateUser(&user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/signin", 301)
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
		fmt.Println("POST!!!")
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		// foundUser := &models.User{}
		var err error
		_, err = app.forum.LogInUser(&user)
		if err != nil {
			app.render(w, r, "signin.page.html", &templateData{
				Error: true,
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

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds
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
