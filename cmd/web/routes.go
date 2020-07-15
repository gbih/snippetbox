package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// APIs
	mux.HandleFunc("/api/v1/snippet", app.apiShowSnippet)
	mux.HandleFunc("/api/v1/test", app.apiTest)
	mux.HandleFunc("/api/v1/home", app.apiHome)

	return mux
}
