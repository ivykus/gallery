package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/context"
	"github.com/ivykus/gallery/errors"
	"github.com/ivykus/gallery/models"
)

type Gallery struct {
	Template struct {
		New  Template
		Edit Template
	}
	GalleryService *models.GalleryService
}

func (g Gallery) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Template.New.Execute(w, r, data)
}

func (g Gallery) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	userId := context.User(r.Context()).Id
	gallery, err := g.GalleryService.Create(data.Title, userId)
	if err != nil {
		g.Template.New.Execute(w, r, data)
		return
	}
	editURL := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editURL, http.StatusFound)
}

func (g Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		fmt.Println("gallery edit", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.Id {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}
	var data struct {
		Title string
		ID    int
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	g.Template.Edit.Execute(w, r, data)
}
