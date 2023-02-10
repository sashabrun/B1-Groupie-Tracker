package main

import (
	"Groupie-Tracker/handlers"
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	http.HandleFunc("/home", handlers.HomeHandler)
	fmt.Println("Listening on http://localhost" + PORT + "/home")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
