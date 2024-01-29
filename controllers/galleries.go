package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ivykus/gallery/context"
	"github.com/ivykus/gallery/errors"
	"github.com/ivykus/gallery/models"
)

type Gallery struct {
	Template struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
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

func (g Gallery) Show(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("gallery show", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var data struct {
		ID     int
		Title  string
		Images []string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title
	for i := 0; i < 15; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		randImgURL := fmt.Sprintf("https://placeimg.com/%d/%d/any", w, h)
		data.Images = append(data.Images, randImgURL)
	}

	g.Template.Show.Execute(w, r, data)
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

func (g Gallery) Update(w http.ResponseWriter, r *http.Request) {
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

	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		fmt.Println("gallery edit", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editURL := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editURL, http.StatusFound)
}

func (g Gallery) Index(w http.ResponseWriter, r *http.Request) {
	type gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []gallery
	}

	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.Id)
	if err != nil {
		fmt.Println("gallery index", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, gal := range galleries {
		data.Galleries = append(data.Galleries, gallery{
			ID:    gal.ID,
			Title: gal.Title,
		})
	}

	g.Template.Index.Execute(w, r, data)
}
