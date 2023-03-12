package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
	_ "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Artists    []Artist
	Categories map[string][]Artist
	Favorites  []Artist
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

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "loading.html", data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if favCookie, err := r.Cookie("Fav"); err != nil {
		data.Favorites = make([]Artist, 0)
		favCookie = &http.Cookie{
			Name:  "Fav",
			Value: EncodeFavCookieValue(data.Favorites),
			//The "Fav" cookie has to never expire
			//to save client's stats in his navigator :
			//their favorite artists.
			Expires: time.Date(2037, 12, 01, 00, 00, 00, 00, time.UTC),
		}
		http.SetCookie(w, favCookie)
	} else {
		DecodeFavCookie(favCookie)
		fmt.Println("Client's fav artists :")
		for _, artist := range data.Favorites {
			fmt.Println(artist.Name)
		}
	}
	_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "artists.gohtml", data)
}
func FavHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "favs.gohtml", data)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	id, _ := strconv.Atoi(parts[2])
	artist := data.Artists[id-1]
	_ = r.ParseForm()
	if strId := r.FormValue("addFav"); strId != "" {
		Id, _ := strconv.Atoi(strId)
		data.Favorites = append(data.Favorites, data.Artists[Id-1])
	}
	_ = tpl.ExecuteTemplate(w, "artist.gohtml", artist)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "error.html", data)
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

	for _, dataArtist := range data.Artists {

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
		for _, genre := range fullartist.Genres {
			data.Categories[genre] = append(data.Categories[genre], dataArtist)
		}
	}
	for category, artists := range data.Categories {
		if len(artists) < 6 {
			delete(data.Categories, category)
		} else {
			fmt.Println(category, ":")
			for _, artist := range artists {
				fmt.Println(artist.Name)
			}
		}
	}
	storeCategories()
}
func storeCategories() {
	CategoriesJSON, _ := json.Marshal(data.Categories)
	if err := os.WriteFile("data/categories.json", CategoriesJSON, 0777); err != nil {
		fmt.Println(err)
	}
}
func GetCategories() {
	data.Categories = make(map[string][]Artist)
	file, _ := os.ReadFile("data/categories.json")
	if len(file) != 0 {
		_ = json.Unmarshal(file, &data.Categories)
	} else {
		ApiCategoryFill()
		_ = json.Unmarshal(file, &data.Categories)
	}
}
