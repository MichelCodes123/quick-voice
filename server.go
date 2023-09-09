package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// var templates = template.Must(template.ParseFiles("templates/index.html", "templates/analytics.html"))
var templates = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, dir string, data any) {
	err := templates.ExecuteTemplate(w, dir, nil)
	//Handling errors associated with ParseFiles method
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func initt() string {
	err:= godotenv.Load()
	if (err != nil){
		log.Fatal("Error loading file")
	}

	str := os.Getenv("PASS")
	dbname := os.Getenv("DBNAME")
	user := os.Getenv("USER")
    port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, str, dbname)
	db, err := sql.Open("postgres", connStr)

	_, e := db.Exec("INSERT into sender VALUES (3, '18 driveOn Road','647-890-1232', 'Joe Smith')")

	db.Close()
	if err != nil {
		log.Fatal(err)
	}
	if e != nil {
		return "something went wrong"
	}
	return ""
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

	http.HandleFunc("/presets", func(w http.ResponseWriter, r *http.Request) {
		initt()
	})

	//Sets up port for listening
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}
