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

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.LoadAndSave)

	mux := pat.New()

	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/original", app.homeOriginal)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	// Static files do not need dynamic middleware
	fileServer := http.FileServer(assetsFileSystem{http.Dir("./ui/static/")})
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// APIs
	mux.Get("/api/v1/snippet", http.HandlerFunc(app.apiShowSnippet))
	mux.Get("/api/v1/test", http.HandlerFunc(app.apiTest))
	mux.Get("/api/v1/home", http.HandlerFunc(app.apiHome))
	// mux.HandleFunc("/api/v1/snippet", app.apiShowSnippet)
	// mux.HandleFunc("/api/v1/test", app.apiTest)
	// mux.HandleFunc("/api/v1/home", app.apiHome)

	// SCS Session test
	// mux.Get("/get", http.HandlerFunc(getHandler))
	// mux.Get("/put", http.HandlerFunc(putHandler))

	// May need to wrap with session.LoadAndSave()
	mux.Get("/get", dynamicMiddleware.ThenFunc(getHandler))
	mux.Get("/put", dynamicMiddleware.ThenFunc(putHandler))

	// return (app.recoverPanic(
	// 	app.logRequest(
	// 		secureHeaders(mux))))

	// Return the alice "standard" middleware chain, followed by servemux
	return standardMiddleware.Then(mux)
}
