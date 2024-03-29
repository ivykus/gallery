package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
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
	gallery, err := g.getGalleryByID(w, r)
	if err != nil {
		return
	}

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}

	var data struct {
		ID     int
		Title  string
		Images []Image
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println("gallery show", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       gallery.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

	g.Template.Show.Execute(w, r, data)
}

func (g Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}

	var data struct {
		ID     int
		Title  string
		Images []Image
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println("gallery edit", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       gallery.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Template.Edit.Execute(w, r, data)
}

func (g Gallery) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r, userMustOwnGallery)
	if err != nil {
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

func (g Gallery) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		fmt.Println("gallery delete", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
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

func (g Gallery) Image(w http.ResponseWriter, r *http.Request) {
	galleryId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	filename := g.getFilename(w, r)

	image, err := g.GalleryService.Image(galleryId, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println("gallery image", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

func (g Gallery) UploadImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("gallery upload image", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("gallery upload image", err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		fmt.Printf("Uploading file: %s\n for gallery: %d\n", fileHeader.Filename, gallery.ID)

		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileError models.FileError
			if errors.As(err, &fileError) {
				msg := fmt.Sprintf("%v has invalid content type or extension. Only jpeg, png and gif are allowed", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			fmt.Println("gallery upload image", err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Gallery) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.getFilename(w, r)

	gallery, err := g.getGalleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		// if errors.Is(err, models.ErrNotFound) {
		//TODO maybe show to user that image was not found
		// }
		fmt.Println("gallery delete image", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Gallery) getFilename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Gallery) getGalleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		fmt.Println("gallery show", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(w, r, gallery); err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.Id {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user %d does not own gallery %d", user.Id, gallery.ID)
	}
	return nil
}
