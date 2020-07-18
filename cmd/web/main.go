package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gbih/snippetbox/pkg/models/postgres"
	_ "github.com/lib/pq"
)

// https://phoenixframework.readme.io/v0.13.1/docs/sessions
// Akin to Phoenix, we can use 2 styles of sessions,
// 1. Client-side cookie-session storage. The analogy in Go here is
// https://github.com/golangcollege/sessions
// 2. Server-side sessions via ETS. The analogy in Go would be:
// https://github.com/alexedwards/scs

type assetsFileSystem struct {
	fs http.FileSystem
}

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *scs.SessionManager
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

var session *scs.SessionManager

func main() {

	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost/snippetbox?sslmode=disable", "Postgres data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session = scs.New()
	session.Store = memstore.New()

	app := &application{
		session:       session,
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgres.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	//fmt.Println("**** templateCache: ", app.templateCache)

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// infoLog.Printf("Server started, http://localhost%v", cfg.Addr)
	// infoLog.Printf("http://localhost%v/original", cfg.Addr)
	// infoLog.Printf("http://localhost%v", cfg.Addr)
	// infoLog.Printf("curl 'http://localhost%v/api/v1/test'", cfg.Addr)
	// infoLog.Printf("curl 'http://localhost%v/api/v1/snippet?id=3'", cfg.Addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

// go run ./cmd/web >> ./logs/info.log 2>> ./logs/error.log
