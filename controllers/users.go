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

	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
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

	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie := http.Cookie{
		Name:     "gallery-session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u User) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
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

	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "gallery-session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "Authenticated user %v\n", user)
}

func (u User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("gallery-session")
	if err != nil {
		fmt.Printf("Error getting cookie: %v\n", err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(tokenCookie.Value)
	if err != nil {
		fmt.Printf("Error getting session: %v\n", err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %v\n", user)
}
