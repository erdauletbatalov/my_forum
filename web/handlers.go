package web

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/pkg/models"
	"forum/pkg/session"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		isSession, user_id := session.IsSession(r)
		var user *models.User
		var err error
		if isSession {
			user, err = app.Forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		}
		sortBy := r.URL.Query().Get("sort")
		posts := []*models.Post{}
		if strings.Compare(sortBy, "likes") == 0 {
			fmt.Println("likes sortion processing")
			posts, err = app.Forum.GetPostsSortedByLikes(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		} else if strings.Compare(sortBy, "date") == 0 {
			fmt.Println("date sortion processing")
			posts, err = app.Forum.GetPostsSortedByDate(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		} else if strings.Compare(sortBy, "tags") == 0 {
			fmt.Println("tag sortion processing")
			tag := r.URL.Query().Get("tag")
			posts, err = app.Forum.GetPostsByTag(user_id, tag)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		} else {
			posts, err = app.Forum.GetPosts(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		}

		tags, err := app.Forum.GetTags()
		if err != nil {
			fmt.Println(err.Error())
			app.serverError(w, err)
			return
		}

		app.render(w, r, "home.page.html", &templateData{
			IsSession: isSession,
			User:      user,
			Posts:     posts,
			Tags:      tags,
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
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Username: r.FormValue("nickname"),
		}
		if !validEmail(user.Email) || user.Password == "" || user.Username == "" {
			app.clientError(w, http.StatusBadRequest)
			return
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

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user_id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || user_id < 1 {
			app.notFound(w)
			return
		}
		isSession, session_user_id := session.IsSession(r)
		user, err := app.Forum.GetUserByID(user_id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.notFound(w)
				return
			}
			app.serverError(w, err)
			return
		}

		filter := r.URL.Query().Get("filter")
		posts := []*models.Post{}
		if strings.Compare(filter, "liked") == 0 {
			posts, err = app.Forum.GetLikedUserPosts(session_user_id, user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		} else if strings.Compare(filter, "posts") == 0 {
			posts, err = app.Forum.GetUserPosts(session_user_id, user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
		}

		app.render(w, r, "profile.page.html", &templateData{
			IsSession: isSession,
			User:      user,
			Posts:     posts,
		})
		return
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
				app.serverError(w, err)
				return
			}
		}
		post, err := app.Forum.GetPostByID(post_id, user_id)
		if err != nil {
			fmt.Println(err.Error())
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		comments, err := app.Forum.GetCommentsByPostID(post_id, user_id)
		if err != nil {
			fmt.Println(err.Error())
			app.serverError(w, err)
			return
		}

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
				app.serverError(w, err)
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
				app.serverError(w, err)
				return
			}
			tagsStr := r.FormValue("tags")
			tagsArr := strings.Split(tagsStr, " ")
			post := models.Post{
				User_id: user.ID,
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Tags:    tagsArr,
			}
			if post.Title == "" || post.Content == "" {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			post.Tags = unique(post.Tags)
			if !validTags(post.Tags) {
				app.render(w, r, "createpost.page.html", &templateData{
					IsSession: isSession,
					IsError:   isError{true, "terrible tags!!!"},
					User:      user,
				})
			}
			if len(post.Tags) > 6 {
				app.render(w, r, "createpost.page.html", &templateData{
					IsSession: isSession,
					IsError:   isError{true, "more than 6 tags are forbidden"},
					User:      user,
				})
			}

			post_id, err := app.Forum.AddPost(&post)
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusSeeOther)
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
				app.serverError(w, err)
				return
			}
			post_id, err := strconv.Atoi(r.FormValue("post_id"))
			if err != nil || post_id < 1 {
				app.serverError(w, err)
				return
			}
			comment := models.Comment{
				User_id: user.ID,
				Post_id: post_id,
				Content: r.FormValue("comment"),
			}
			if comment.Content == "" {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			err = app.Forum.AddComment(&comment)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
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
		isSession, user_id := session.IsSession(r)
		if isSession {
			user, err := app.Forum.GetUserByID(user_id)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
			post_id, err := strconv.Atoi(r.FormValue("post_id"))
			if err != nil || post_id < 1 {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			comment_id, err := strconv.Atoi(r.FormValue("comment_id"))
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			vote_type, err := strconv.Atoi(r.FormValue("vote_type"))
			if err != nil {
				fmt.Println(err.Error())
				app.clientError(w, http.StatusBadRequest)
				return
			} else if vote_type != 1 && vote_type != -1 {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			vote := models.Vote{
				User_id:    user.ID,
				Post_id:    post_id,
				Comment_id: comment_id,
				Vote_type:  vote_type,
			}
			vote_type, err = app.Forum.GetVoteType(&vote)
			if err != nil {
				fmt.Println(err.Error())
				app.serverError(w, err)
				return
			}
			switch vote_type {
			case 0: // there is no votes yet
				err = app.Forum.AddVote(&vote)
				if err != nil {
					fmt.Println(err.Error())
					app.serverError(w, err)
					return
				}
			case 1: // there is Like
				if vote.Vote_type == 1 { // trying to like when there already a like
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
				} else { // like when there is no like
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
					err = app.Forum.AddVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
				}
			case -1: // there is Dislike
				if vote.Vote_type == -1 { // trying to dislike when there already a dislike
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
				} else { // dislike when there is no dislike
					err = app.Forum.DeleteVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
					err = app.Forum.AddVote(&vote)
					if err != nil {
						fmt.Println(err.Error())
						app.serverError(w, err)
						return
					}
				}
			default: // Post rate
				fmt.Println(err.Error())
				app.serverError(w, err)
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
