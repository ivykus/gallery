package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
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

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
