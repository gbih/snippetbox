package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gbih/snippetbox/pkg/forms"
	"github.com/gbih/snippetbox/pkg/models"
)

func (nfs assetsFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

//----------

// API Usage: curl -i localhost:4000/api/v1/test

func (app *application) apiTest(w http.ResponseWriter, r *http.Request) {
	data := []struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}{
		{"George", "日本"},
		{"Izumi", "越後湯沢"},
	}

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

//----------

// Example to show explicit code flow

func (app *application) homeOriginal(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/home.pageOriginal.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	// Create a FuncMap with which to register the function.
	// "humanDate" is the function name created in templates.go,
	//  to be called in the template text.
	funcMap := template.FuncMap{
		"humanDate": humanDate,
	}

	// Create a template, add the function map, and parse the text.
	// In order to use .Funcs(), we have to call template.New()
	// The template.FuncMap must be registered with the template set before
	// calling the ParseFiles() method. So, we have to use template.New() to
	// create an empty template set, use the Funcs() method to register the
	// template.FuncMap, and then parse the file as normal.
	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the template, which is essentially a function.
	// As soon as we begin adding dynamic behavior to our HTML templates
	// there’s a risk of encountering runtime errors. We want to avoid
	// the use-case of the template compiling ok, but throwing a run-time error.
	// Essentially, To fix this we need to make the template render a two-stage
	// process. First, we should make a trial-render by writing the template
	// into a buffer. If this fails, we can respond to the user with an error
	// message. But if it works, we can then write the contents of the buffer
	// to our http.ResponseWriter.

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper
	// and then return.
	err = ts.Execute(buf, &templateData{Snippets: s})
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)

}

//----------

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// New code
	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})

	/*
		Note:
		1. 	Old template parsing and execution code refactored out, moved to
		templates.go in this function:
		newTemplateCache(dir string) (map[string]*template.Template, error)

		2. This cache is initialized in the main() function and made
		available to handlers as a dependency via the application struct.

			type application struct {
				templateCache map[string]*template.Template
			}

		and added as to the application dependencies
			app := &application{
				templateCache: templateCache,
			}

		So now we have an in-memory cache of the relevant template set for
		each of our pages, and our handlers have access to this cache via
		the application struct.

		Now, create a helper render method so that we can easily render
		the templates from the cache.
		File: cmd/web/helpers.go

		func (app *application) render(
				w http.ResponseWriter,
				r *http.Request,
				name string,
				td *templateData) { ... }

		Now, we can access these cached templates via the render helper function:
		Use the new render helper.
		app.render(w, r, "home.page.tmpl", &templateData{
			Snippets: s,
		})
	*/
	/*
		Old code:
		data := &templateData{Snippets: s}

		files := []string{
			"./ui/html/home.page.html",
			"./ui/html/base.layout.html",
			"./ui/html/footer.partial.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err)
		}
	*/
}

// API Usage : curl -i localhost:4000/api/v1/home

func (app *application) apiHome(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	js, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json; charset=UTF-8")
	w.Write(js)
}

//----------

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Manual template parsing and execution code refactored out

	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})

}

// API Usage: curl -i 'http://localhost:4000/api/v1/snippet?id=1'

func (app *application) apiShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	js, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "%v", s)
	w.Header().Set("Content-Type", "application-json; charset=utf-8")
	w.Write(js)
}

//----------

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// ParseForm populates r.Form and r.PostForm.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	form.PermittedValues("items", "foo", "bar", "baz")

	// Multiple-Value Fields
	// Strictly speaking, the r.PostForm.Get() method that we’ve used above only returns the first value for a specific form field. This means you can’t use it with form fields which potentially send multiple values, such as a group of checkboxes.
	// In this case you’ll need to work with the r.PostForm map directly. The underlying type of the r.PostForm map is url.Values, which in turn has the underlying type map[string][]string. So, for fields with multiple values you can loop over the underlying map to access them like so:

	// DEBUG
	// for i, item := range r.PostForm["items"] {
	// 	fmt.Fprintf(w, "%d: Item %s\n", i, item)
	// }
	const foo = "foo"
	const bar = "bar"
	const baz = "baz"
	const blank = "---"

	// the r.PostForm.Get() method only returns the first value for a specific form field. This means you can’t use it with form fields which potentially send multiple values, such as a group of checkboxes.
	// Here we need to work with the r.PostForm map directly. The underlying type of the r.PostForm map is url.Values, which in turn has the underlying type map[string][]string.

	// Checkboxes (many-in-a-set) validation
	// https://www.alexedwards.net/blog/validation-snippets-for-go#checkboxes-many-in-set
	// GB Map is the correct data structure here, as it enables using a
	// key-value pair that we can then iterate through and check
	//set := map[string]bool{"foo": true, "baz": true, "bar": true}
	// var foo string
	// var baz string
	// var bar string

	for _, item := range r.PostForm["items"] {
		// have to manually add to form data structure so we can
		// more easily check and manipulate it
		form.Add(item, item)
	}

	if form.Get(foo) != foo {
		form.Set(foo, blank)
	}
	if form.Get(bar) != bar {
		form.Set(bar, blank)
	}
	if form.Get(baz) != baz {
		form.Set(baz, blank)
	}

	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := app.snippets.Insert(
		form.Get("title"),
		form.Get("content"),
		form.Get("expires"),
		form.Get("foo"),
		form.Get("bar"),
		form.Get("baz"),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	//
	// Cookie style:
	// app.session.Put(r, "flash", "Snippet successfully created!")
	//
	// Server-side style:Store a new key and value in the session data.
	session.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
