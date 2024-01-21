package controllers

import (
	"fmt"
	"net/http"

	"github.com/ivykus/gallery/models"
)

type User struct {
	Templates struct {
		New Template
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
