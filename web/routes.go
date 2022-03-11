package web

import (
	"net/http"
)

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signup", app.signup)
	mux.HandleFunc("/signin", app.signin)
	mux.HandleFunc("/signout", app.signout)
	mux.HandleFunc("/user", app.profile)
	mux.HandleFunc("/post/create", app.createPost)
	mux.HandleFunc("/post", app.showPost)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
