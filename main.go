package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Home page</h1>")
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Contacts page</h1>")
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contacts":
		contactsHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", router)
}
