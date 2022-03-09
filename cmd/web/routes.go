package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signup", app.signup)
	mux.HandleFunc("/signin", app.signin)
	mux.HandleFunc("/signout", app.signout)
	mux.HandleFunc("/user/profile", app.profile)
	mux.HandleFunc("/post/create", app.createPost)
	mux.HandleFunc("/post", app.showPost)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
