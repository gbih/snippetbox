package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/gbih/snippetbox/pkg/alice"
)

func (app *application) routes() http.Handler {

	// Use alice to create our "standard" middleware chain,
	// used for every request our app receives. This is just our
	// middleware arranged in a list or slice data structure.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/original", app.homeOriginal)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(assetsFileSystem{http.Dir("./ui/static/")})
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// APIs
	mux.Get("/api/v1/snippet", http.HandlerFunc(app.apiShowSnippet))
	mux.Get("/api/v1/test", http.HandlerFunc(app.apiTest))
	mux.Get("/api/v1/home", http.HandlerFunc(app.apiHome))
	// mux.HandleFunc("/api/v1/snippet", app.apiShowSnippet)
	// mux.HandleFunc("/api/v1/test", app.apiTest)
	// mux.HandleFunc("/api/v1/home", app.apiHome)

	// return (app.recoverPanic(
	// 	app.logRequest(
	// 		secureHeaders(mux))))

	// Return the alice "standard" middleware chain, followed by servemux
	return standardMiddleware.Then(mux)
}
