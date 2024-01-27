package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/ivykus/gallery/controllers"
	"github.com/ivykus/gallery/migrations"
	"github.com/ivykus/gallery/models"
	"github.com/ivykus/gallery/templates"
	"github.com/ivykus/gallery/views"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Adress string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.DefaultPostgresConfig()
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: read values from env

	cfg.CSRF.Key = "FnsdflDSD9SDg82nlz00guu23xvjsDdD"
	cfg.CSRF.Secure = true

	cfg.Server.Adress = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// setup database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFs(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	userService := &models.UserService{DB: db}
	sessionService := &models.SessionService{DB: db}
	pwResetService := &models.PasswordResetService{DB: db}
	emailService := models.NewEmailService(cfg.SMTP)

	// setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// setup controllers
	usersC := controllers.User{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml"))

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
	fmt.Println("Server is running on port ", cfg.Server.Adress)
	err = http.ListenAndServe(cfg.Server.Adress, r)
	if err != nil {
		panic(err)
	}
}
