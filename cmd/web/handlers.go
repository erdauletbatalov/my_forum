package main

import (
	"fmt"
	"net/http"

	"github.com/erdauletbatalov/forum.git/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "home.page.html", &templateData{})
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
