package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// GB: Good example of mutating a field of an instance of a struct.
// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. We're not using the *http.Request parameter at the moment,
// but we will do later in the book.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()

	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	td.Flash = session.PopString(r.Context(), "flash")

	return td
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Need to make the template render a two-stage process. First, we should
	// make a ‘trial’ render by writing the template into a buffer. If this
	// fails, we can respond to the user with an error message. But if it works,
	// we can then write the contents of the buffer to our http.ResponseWriter.

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.
	//
	// Execute the template set, passing the dynamic data with the current
	// year injected.
	//err := ts.Execute(buf, td)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)

}
