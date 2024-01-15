package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/views"
)

func ExecuteTemplate(w http.ResponseWriter, path string) {
	viewTmpl, err := views.Parse(path)
	if err != nil {
		log.Printf("%q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	viewTmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "home.gohtml")
	ExecuteTemplate(w, path)
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "contacts.gohtml")
	ExecuteTemplate(w, path)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Println("id = ", id)
	fmt.Fprintf(w, "<p>Id: %s", id)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contacts", contactsHandler)
	r.Route("/user", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", userHandler)
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
