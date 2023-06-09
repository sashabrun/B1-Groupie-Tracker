package handlers

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
	_ "log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Artists    []Artist
	Categories map[string][]Artist
	Input      input
	FavIndexs  []int
	Likes      []int
	foundcount int
}
type input struct {
	text      string
	creaDate  string
	nbMembers string
}
type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type Artist struct {
	Id            int      `json:"id"`
	Image         string   `json:"image"`
	Category      []string `json:"category"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Locations     string   `json:"locations"`
	ConcertDates  string   `json:"concertDates"`
	RelationsLink string   `json:"relations"`
	Relations     Relations
	MostListened  string `json:"mostListened"`
	Isliked       bool   `json:"isliked"`
}

var tpl = template.Must(template.New("").Funcs(template.FuncMap{
	"ArtistNameContainsInput": ArtistNameContainsInput,
	"DisplayLocationLink":     DisplayLocationLink,
	"GetArtistLikes":          GetArtistLikes,
}).ParseGlob("web/templates/*"))
var data Data

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "loading.html", data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	CheckFavCookie(r)
	_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	CheckFavCookie(r)
	data.foundcount = 0
	data.Input = input{}
	_ = r.ParseForm()
	if textInput := r.FormValue("research-text"); textInput != "" {
		data.Input.text = textInput
	}
	if dateInput := r.FormValue("range"); dateInput != "" {
		data.Input.creaDate = dateInput
	}
	if nbMembers := r.FormValue("nb-members"); nbMembers != "" {
		data.Input.nbMembers = nbMembers
	}
	fmt.Println("Text Input:", data.Input.text)
	fmt.Println("NbMembers Input:", data.Input.nbMembers)
	fmt.Println("CreaDate Input:", data.Input.creaDate)
	_ = tpl.ExecuteTemplate(w, "artists.gohtml", data)
}
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	CheckFavCookie(r)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	id, _ := strconv.Atoi(parts[2])
	id--
	if r.Method == "POST" {
		_ = r.ParseForm()
		if r.FormValue("addFav") != "" {
			if data.Artists[id].Isliked {
				data.Likes[id]--
				fmt.Println(data.Artists[id].Name, "has now", data.Likes[id], "likes")
				SaveLikes()
				removeFav(id)
				fmt.Println(data.Artists[id].Name, "Is no more in ur list")
			} else {
				data.Likes[id]++
				fmt.Println(data.Artists[id].Name, "has now", data.Likes[id], "likes")
				SaveLikes()
				data.FavIndexs = append(data.FavIndexs, id)
				fmt.Println(data.Artists[id].Name, "Is now in ur list")
			}
			data.Artists[id].Isliked = !data.Artists[id].Isliked
			UpdateFavCookie(w)
		}
	}
	_ = tpl.ExecuteTemplate(w, "artist.gohtml", data.Artists[id])
}

func MyListHandler(w http.ResponseWriter, r *http.Request) {
	CheckFavCookie(r)
	_ = tpl.ExecuteTemplate(w, "mylist.gohtml", data)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "error.html", data)
}
func MostLikedHandler(w http.ResponseWriter, r *http.Request) {
	CheckFavCookie(r)
	_ = tpl.ExecuteTemplate(w, "mostliked.gohtml", data)
}

func ApiCategoryFill() {
	data.Categories = make(map[string][]Artist)
	// Créer une configuration client credentials pour l'authentification OAuth2
	config := &clientcredentials.Config{
		ClientID:     "216d0c4dde684593ba59c143326593aa",
		ClientSecret: "8539b04a599946cbb4840c96058544f9",
		TokenURL:     spotify.TokenURL,
	}

	// Obtenir un client authentifié avec la configuration client credentials
	token, err := config.Token(context.Background())
	if err != nil {
		fmt.Println("Erreur d'obtention de jeton d'accès:", err)
		os.Exit(1)
	}
	client := spotify.Authenticator{}.NewClient(token)

	for i, dataArtist := range data.Artists {
		// Rechercher l'artiste sur Spotify
		results, err := client.Search(dataArtist.Name, spotify.SearchTypeArtist)
		if err != nil {
			fmt.Println("Erreur de recherche d'artiste:", err)
			os.Exit(1)
		}

		// Vérifier si des artistes ont été trouvés
		if len(results.Artists.Artists) == 0 {
			fmt.Println("Aucun artiste trouvé pour", dataArtist.Name)
			os.Exit(1)
		}

		// Sélectionner le premier artiste trouvé
		artist := results.Artists.Artists[0]

		// Récupérer les catégories de musique de l'artiste
		fullartist, err := client.GetArtist(artist.ID)
		if err != nil {
			fmt.Println("Erreur de récupération des catégories de musique de l'artiste:", err)
			os.Exit(1)
		}
		topTracks, err := client.GetArtistsTopTracks(artist.ID, "US")
		if err != nil {
			fmt.Println("Error getting toptracks")
		}
		data.Artists[i].MostListened = string(topTracks[0].ID)
		for _, genre := range fullartist.Genres {
			data.Categories[genre] = append(data.Categories[genre], dataArtist)
		}
	}

	for category, artists := range data.Categories {
		if len(artists) < 6 {
			delete(data.Categories, category)
		} else {
			fmt.Println(category, ":")
			for i, _ := range artists {
				fmt.Println(data.Artists[i].Name)
			}
		}
	}
	storeCategories()
}
