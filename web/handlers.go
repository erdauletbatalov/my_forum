package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"forum/pkg/models"
	"forum/pkg/session"

	uuid "github.com/satori/go.uuid"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		posts, err := app.Forum.GetAllPosts()
		if err != nil {
			fmt.Println(err.Error())
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.Forum.GetUserByID(user_id)
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

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("POST!!!")
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Username: r.FormValue("nickname"),
		}
		err := app.Forum.AddUser(&user)
		if err != nil {
			switch err.Error() {
			case "UNIQUE constraint failed: user.email":
				app.render(w, r, "signup.page.html", &templateData{
					IsError: isError{true, "this email is already in use"},
				})
				return
			case "UNIQUE constraint failed: user.username":
				fmt.Println("rendering username already in use")
				app.render(w, r, "signup.page.html", &templateData{
					IsError: isError{true, "this username is already in use"},
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

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		info := r.FormValue("email")
		password := r.FormValue("password")

		err := app.Forum.PasswordCompare(info, password)
		if err != nil {
			fmt.Println(err.Error())
			app.render(w, r, "signin.page.html", &templateData{
				IsError: isError{true, "incorrect email or password"},
			})
			return
		}

		u, _ := app.Forum.GetUserInfo(info)

		sessionToken := uuid.NewV4().String()
		expiresAt := time.Now().Add(120 * time.Second)

		session.Sessions[sessionToken] = session.Session{
			ID:     u.ID,
			Expiry: expiresAt,
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		http.Redirect(w, r, fmt.Sprintf("/user?id=%v", session.Sessions[sessionToken].ID), http.StatusSeeOther)
	case http.MethodGet:
		app.render(w, r, "signin.page.html", &templateData{})
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user_id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || user_id < 1 {
			app.notFound(w)
			return
		}
		isSession, _ := session.IsSession(r)
		user, err := app.Forum.GetUserByID(user_id)
		if err != nil {
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		app.render(w, r, "profile.page.html", &templateData{
			IsSession: isSession,
			User:      user,
		})
		return
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) signout(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) showPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		post_id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || post_id < 1 {
			app.notFound(w)
			return
		}
		isSession, user_id := session.IsSession(r)
		var user *models.User
		if isSession {
			user, err = app.Forum.GetUserByID(user_id)
			if err != nil {
				app.clientError(w, http.StatusInternalServerError)
				return
			}
		}
		post, err := app.Forum.GetPostByID(post_id)
		if err != nil {
			fmt.Println("getPostByID fail")
			fmt.Println(err.Error())
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		fmt.Println("getPostByID success")

		comments, err := app.Forum.GetCommentsByPostID(post_id, user_id)
		if err != nil {
			fmt.Println("GetCommentsByPostID fail")
			fmt.Println(err.Error())
			app.clientError(w, http.StatusInternalServerError)
			return
		}

		fmt.Println("GetCommentsByPostID success")
		// Используем помощника render() для отображения шаблона.
		app.render(w, r, "post.page.html", &templateData{
			IsSession: isSession,
			User:      user,
			Post:      post,
			Comments:  comments,
		})
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		isSession, user_id := session.IsSession(r)

		if isSession {
			user, err := app.Forum.GetUserByID(user_id)
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
			user, err := app.Forum.GetUserByID(user_id)
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
			id, err := app.Forum.AddPost(&post)
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

func (app *Application) createComment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.Forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			post_id, err := strconv.Atoi(r.FormValue("post_id"))
			if err != nil || post_id < 1 {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			fmt.Printf("MethodPost post_id = %v\n", post_id)
			comment := models.Comment{
				User_id: user.ID,
				Post_id: post_id,
				Content: r.FormValue("comment"),
			}
			if comment.Content == "" {
				fmt.Println("Empty comment error")
				app.clientError(w, http.StatusBadRequest)
				return
			}
			err = app.Forum.AddComment(&comment)
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusSeeOther)
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *Application) rate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("entering rate handler")
		isSession, user_id := session.IsSession(r)
		if isSession {
			fmt.Println("entering rate handler, you are loged in")

			user, err := app.Forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			post_id, err := strconv.Atoi(r.FormValue("post_id"))
			if err != nil || post_id < 1 {
				fmt.Println("post_id err")
				app.clientError(w, http.StatusBadRequest)
				return
			}
			comment_id, err := strconv.Atoi(r.FormValue("comment_id"))
			if err != nil {
				fmt.Println("comment_id err")
				app.clientError(w, http.StatusBadRequest)
				return
			}
			vote_type, err := strconv.Atoi(r.FormValue("vote_type"))

			if err != nil {
				fmt.Println("vote_type err")
				fmt.Println(err.Error())
				app.clientError(w, http.StatusBadRequest)
				return
			} else if vote_type != 1 && vote_type != -1 {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			vote_obj, err := strconv.Atoi(r.FormValue("vote_obj"))
			if err != nil {
				fmt.Println("vote_obj err")
				app.clientError(w, http.StatusBadRequest)
				return
			}
			fmt.Printf("MethodPost post_id = %v\n", post_id)
			vote := models.Vote{
				User_id:    user.ID,
				Post_id:    post_id,
				Comment_id: comment_id,
				Vote_type:  vote_type,
				Vote_obj:   vote_obj,
			}
			fmt.Println("entering GetVoteType")
			vote_type, err = app.Forum.GetVoteType(&vote)
			if err != nil {
				fmt.Println("GetVoteType error")
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
				return
			}
			fmt.Println("GetVoteType success, vote type:")
			fmt.Println(vote_type)
			switch vote_type {
			case 0: // there is no votes yet
				fmt.Println("entering AddVote")
				err = app.Forum.AddVote(&vote)
				if err != nil {
					fmt.Println(err.Error())
					app.clientError(w, http.StatusInternalServerError)
					return
				}
			case 1: // there is Like
				if vote.Vote_type == 1 { // trying to like when there already a like
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
				} else { // like when there is no like
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
					err = app.Forum.AddVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
				}
			case -1: // there is Dislike
				if vote.Vote_type == -1 { // trying to dislike when there already a dislike
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
				} else { // dislike when there is no dislike
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
					err = app.Forum.AddVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.clientError(w, http.StatusInternalServerError)
						return
					}
				}
			default: // Post rate
				fmt.Println(err.Error())
				app.clientError(w, http.StatusInternalServerError)
			}

			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
