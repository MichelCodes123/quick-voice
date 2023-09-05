package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/analytics.html"))

func renderTemplate(w http.ResponseWriter, dir string, data any) {
	err := templates.ExecuteTemplate(w, dir, nil)
	//Handling errors associated with ParseFiles method
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	//Registering route handlers
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//Setup as the root directory for the app, rejects any request that is not defined in the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		if r.URL.Path != "/" {
			http.Error(w, "404 not found", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "index.html", nil)
	})

	//http package provides the HandleFunc method, which accepts a path and a handler containing response and request information.
	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "analytics.html", nil)
	})

	//Sets up port for listening
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}
