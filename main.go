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
	http.HandleFunc("/fav-artists", handlers.FavHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/", handlers.ErrorHandler)

	handlers.FillData()
	handlers.GetCategories()
	fmt.Println(handlers.DisplayLocationLink(24))

	fmt.Println("Listening on http://localhost" + PORT + "/loading")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
