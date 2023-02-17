package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	Artists []Artist
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var tpl = template.Must(template.ParseGlob("web/templates/*"))
var data Data

func FillData() {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err == nil {
		_ = json.NewDecoder(res.Body).Decode(&data.Artists)
	} else {
		log.Fatal(err)
	}
	defer res.Body.Close()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	FillData()
	_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
}
