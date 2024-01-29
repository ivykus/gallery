package controllers

import (
	"fmt"
	"net/http"

	"github.com/ivykus/gallery/context"
	"github.com/ivykus/gallery/models"
)

type Gallery struct {
	Template struct {
		New Template
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
