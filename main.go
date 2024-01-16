package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/controllers"
	"github.com/ivykus/gallery/templates"
	"github.com/ivykus/gallery/views"
)

func main() {
	r := chi.NewRouter()
	fs := templates.FS

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(fs, "home.gohtml"))))

	r.Get("/contacts", controllers.StaticHandler(
		views.Must(views.ParseFS(fs, "contacts.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
