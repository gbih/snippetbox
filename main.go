package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL.Path:", r.URL.Path)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet"))
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("curl localhost:4000/ -i")
	log.Println("curl localhost:4000/snippet -i")
	log.Println("curl localhost:4000/snippet/create -i")
	log.Println("curl localhost:4000/missing-i")
	log.Println("curl -i -X POST http://localhost:4000/snippet/create")
	log.Println("curl -i -X DELETE http://localhost:4000/snippet/create")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
