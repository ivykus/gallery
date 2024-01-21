package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/controllers"
	"github.com/ivykus/gallery/models"
	"github.com/ivykus/gallery/templates"
	"github.com/ivykus/gallery/views"
)

func main() {
	r := chi.NewRouter()
	fs := templates.FS

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(fs, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contacts", controllers.StaticHandler(
		views.Must(views.ParseFS(fs, "contacts.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FaqHandler(
		views.Must(views.ParseFS(fs, "faq.gohtml", "tailwind.gohtml"))))

	cfg := models.DefaultPostgresConfig()

	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	UserService := models.UserService{DB: db}

	usersC := controllers.User{
		UserService: &UserService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(fs, "signup.gohtml",
		"tailwind.gohtml"))

	usersC.Templates.SignIn = views.Must(views.ParseFS(fs, "signin.gohtml",
		"tailwind.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Server is running on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
