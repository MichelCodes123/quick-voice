package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//Registering route handlers
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/",fs)

	//Sets up port for listening
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
