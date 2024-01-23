package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
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

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return `<input type="hidden" name="csrf" />`
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("Error cloning template: %q", err)
		http.Error(w, "Error while rendering page", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
