package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
		posts, err := app.forum.GetAllPosts()
		if err != nil {
			fmt.Println(err.Error())
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			app.render(w, r, "home.page.html", &templateData{
				IsSession: isSession, //
				User:      user,
				Posts:     posts,
			})
			return
		}
		app.render(w, r, "home.page.html", &templateData{
			IsSession: isSession, //
			Posts:     posts,
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
		err := app.forum.AddUser(&user)
		if err != nil {
			switch err.Error() {
			case "UNIQUE constraint failed: user.email":
				app.render(w, r, "signup.page.html", &templateData{
					IsError: isError{true, "this email is already in use"},
				})
				return
			case "UNIQUE constraint failed: user.nickname":
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
		err := app.forum.LogInUser(&user)
		if err != nil {
			fmt.Println(err.Error())
			app.render(w, r, "signin.page.html", &templateData{
				IsError: isError{true, "incorrect email or password"},
			})
			return
		}

		u, _ := app.forum.GetUserByEmail(user.Email)

		// Create a new random session token
		// we use the "github.com/google/uuid" library to generate UUIDs
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		// Set the token in the session map, along with the session information
		session.Sessions[sessionToken] = session.Session{
			ID:     u.ID,
			Expiry: expiresAt,
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
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
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.forum.GetUserByID(user_id)
			if err != nil {
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			app.render(w, r, "profile.page.html", &templateData{
				IsSession: isSession,
				User:      user,
			})
			return
		} else {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
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

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		isSession, user_id := session.IsSession(r)
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}
		user, err := app.forum.GetUserByID(user_id)
		if err != nil {
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		post, err := app.forum.GetPostByID(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		// Используем помощника render() для отображения шаблона.
		app.render(w, r, "post.page.html", &templateData{
			IsSession: isSession,
			Post:      post,
			User:      user,
		})
	default:
	}
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		isSession, user_id := session.IsSession(r)

		if isSession {
			user, err := app.forum.GetUserByID(user_id)
			if err != nil {
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			app.render(w, r, "createpost.page.html", &templateData{
				IsSession: isSession,
				User:      user,
			})
			return
		} else {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
	case http.MethodPost:
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			post := models.Post{
				User_id: user.ID,
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
			}
			id, err := app.forum.AddPost(&post)
			if err != nil {
				fmt.Println(err.Error())
				app.render(w, r, "createpost.page.html", &templateData{
					IsSession: isSession,
					IsError:   isError{true, err.Error()},
					User:      user,
				})
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
		}
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
