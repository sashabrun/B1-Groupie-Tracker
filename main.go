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
	http.HandleFunc("/loading", handlers.LoadingHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/artists", handlers.ArtistsHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/mostliked", handlers.MostLikedHandler)
	http.HandleFunc("/mylist", handlers.MyListHandler)
	http.HandleFunc("/", handlers.ErrorHandler)

	handlers.FillData()
	handlers.CategoryFill()

	fmt.Println("Listening on http://localhost" + PORT + "/loading")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
