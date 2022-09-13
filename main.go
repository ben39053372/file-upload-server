package main

import (
	"html/template"
	"log"
	"net/http"
)

var indexPage = template.Must(template.ParseFiles("public/index.html"))

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			indexPage.ExecuteTemplate(w, "index.html", nil)
		} else if r.Method == "POST" {
			uploadFile(w, r)
		}
	})

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getFile(w, r)
		}
	})
	log.Default().Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
