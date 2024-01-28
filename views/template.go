package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ivykus/gallery/context"
	"github.com/ivykus/gallery/models"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

type Template struct {
	htmlTpl *template.Template
}

type public interface {
	Public() string
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return []string{
				"Error1",
				"Error2",
				"Error3",
			}
		},
	})
	tmpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing error %w", err)
	}

	return Template{
		htmlTpl: tmpl,
	}, nil
}

func Parse(path string) (Template, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return Template{}, fmt.Errorf("parsing error %w", err)
	}

	return Template{
		htmlTpl: tmpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("Error cloning template: %q", err)
		http.Error(w, "Error while rendering page", http.StatusInternalServerError)
		return
	}
	errMsgs := ErrorMessages(errs...)

	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errMsgs
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrorMessages(errs ...error) []string {
	var errMsgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			errMsgs = append(errMsgs, pubErr.Public())
		} else {
			fmt.Println(err)
			errMsgs = append(errMsgs, "Something went wrong...")
		}
	}
	return errMsgs
}
