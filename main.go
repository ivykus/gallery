package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/controllers"
	"github.com/ivykus/gallery/views"
)

func main() {
	r := chi.NewRouter()

	tpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "contacts.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/contacts", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
