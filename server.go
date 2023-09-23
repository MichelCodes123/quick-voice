package main

import (
	"database/sql"
	"encoding/json"
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
type sdr struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Number  string `json:"number"`
	Email string 	`json:"email"`
	Name    string `json:"name"`
}

func initt(w http.ResponseWriter) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	defer db.Close()
	rows, e := db.Query(`SELECT * FROM sender`)
	if e != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	var data sdr
	a := make([]sdr, 0)

	for rows.Next() {
		read_err := rows.Scan(&data.Id, &data.Address, &data.Number, &data.Email, &data.Name)
		//Error handling for issue with database reads
		if read_err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		defer rows.Close()
		a = append(a, data)

	}

	str, re := json.Marshal(a)
	if re != nil {
		fmt.Print(re)
	}

	//Setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writerr := w.Write(str)

	//Error handling for issues with writing the response
	if writerr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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

	//make sure this route only accepts get requests
	http.HandleFunc("/presets", func(w http.ResponseWriter, r *http.Request) {
		initt(w)
	})

	type receiver struct{

	}

	type invoice struct{

	}
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {


		r.ParseForm()

		for key,value := range r.Form{
			fmt.Printf("%s = %s\n", key,value)
		}
	

		switch r.Method{
		case "POST":
				renderTemplate(w, "printout.html", nil)
		default:
			http.Error(w, "404 not found", http.StatusInternalServerError)
		}
	
	})

	//Sets up port for listening
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}
