package controllers

import (
	"fmt"
	"net/http"

	"github.com/ivykus/gallery/models"
)

type User struct {
	Templates struct {
		New    Template
		SignIn Template
	}

	UserService *models.UserService
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	user, err := u.UserService.CreateUser(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Created user %v\n", user)
}

func (u User) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, data)
}

func (u User) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "gallery-email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "Authenticated user %v\n", user)
}

func (u User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("gallery-email")
	if err != nil {
		fmt.Printf("Error getting cookie: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Current user: %v\n", cookie.Value)
}
