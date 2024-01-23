package controllers

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "gallery-session"
)

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}

func readCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		fmt.Printf("Error getting cookie: %v\n", err)
		return "", err
	}
	return cookie.Value, nil
}
