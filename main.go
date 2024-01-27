package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/ivykus/gallery/controllers"
	"github.com/ivykus/gallery/migrations"
	"github.com/ivykus/gallery/models"
	"github.com/ivykus/gallery/templates"
	"github.com/ivykus/gallery/views"
)

func main() {
	// setup database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFs(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	UserService := models.UserService{DB: db}
	SessionService := models.SessionService{DB: db}

	// setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &SessionService,
	}

	csrfKey := "FnsdflDSD9SDg82nlz00guu23xvjsDdD"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		//TODO: fix this before deployment
		csrf.Secure(false),
	)

	// setup controllers
	usersC := controllers.User{
		UserService:    &UserService,
		SessionService: &SessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml"))

	// setup routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml"))))
	r.Get("/contacts", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contacts.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FaqHandler(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml"))))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// start the server
	fmt.Println("Server is running on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
