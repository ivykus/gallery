package main

import (
	"os"
	"text/template"
)

type User struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{"John"}
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
