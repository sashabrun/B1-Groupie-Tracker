package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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
	apiRes, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err == nil {
		_ = json.NewDecoder(apiRes.Body).Decode(&data.Artists)
	}
	defer apiRes.Body.Close()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
}
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	id, _ := strconv.Atoi(parts[2])
	artist := data.Artists[id-1]
	_ = tpl.ExecuteTemplate(w, "artist.gohtml", artist)
}
