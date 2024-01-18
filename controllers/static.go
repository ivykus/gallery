package controllers

import (
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FaqHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, FaqData)
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
