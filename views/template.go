package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
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
	tmpl, err := template.ParseFS(fs, patterns...)
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

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type FaqResponse struct {
	Question string
	Answer   string
}

var FaqData = []FaqResponse{
	{
		Question: "lorem ipsum dolor sit amet?",
		Answer: `labore et dolore magna aliqua ut enim ad minim 
veniam quis nostrud exercitation
ullamco laboris nisi ut aliquip `,
	},
	{
		Question: "voluptate velit esse cillum ",
		Answer:   "dolore eu fugiat nulla pariatur voluptate velit esse cillum",
	},
	{
		Question: "Excepteur sint occaecat cupidatat non proident?",
		Answer:   "sunt in culpa qui officia deserunt mollit anim id est laborum",
	},
	{
		Question: "sunt in culpa qui officia deserunt?",
		Answer: `mollit anim id est laborum sunt in culpa qui officia deserunt
mollit anim id est laborum sunt in culpa qui officia deserunt`,
	},
}
