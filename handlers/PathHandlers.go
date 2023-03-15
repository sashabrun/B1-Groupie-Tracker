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
)

type Data struct {
	Artists    []Artist
	Categories map[string][]Artist
	Input      input
	FavIndexs  []int
	Likes      []int
}
type input struct {
	text string
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
	if favCookie, err := r.Cookie("Fav"); err == nil {
		DecodeFavCookie(favCookie)
		fmt.Println("Client's fav artists :")
		for _, artist := range data.Artists {
			if artist.Isliked {
				fmt.Println(artist.Name)
			}
		}
	} else {
		fmt.Println("No \"Fav\" cookie yet")
	}
	_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if favCookie, err := r.Cookie("Fav"); err == nil {
		DecodeFavCookie(favCookie)
		fmt.Println("Client's fav artists :")
		for _, artist := range data.Artists {
			if artist.Isliked {
				fmt.Println(artist.Name)
			}
		}
	} else {
		fmt.Println("No \"Fav\" cookie yet")
	}
	data.Input.text = ""
	_ = r.ParseForm()
	if textInput := r.FormValue("research-text"); textInput != "" {
		data.Input.text = textInput
	}
	if dateInput := r.FormValue("range"); dateInput != "" {
		fmt.Println(dateInput)
	}
	_ = tpl.ExecuteTemplate(w, "artists.gohtml", data)
}
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
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
	_ = tpl.ExecuteTemplate(w, "mylist.gohtml", data)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "error.html", data)
}
func MostLikedHandler(w http.ResponseWriter, r *http.Request) {
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
func storeArtists() {
	ArtistsJSON, _ := json.Marshal(data.Artists)
	if err := os.WriteFile("data/artists.json", ArtistsJSON, 0777); err != nil {
		fmt.Println(err)
	}
}
func GetCategories() {
	data.Categories = make(map[string][]Artist)
	file, _ := os.ReadFile("data/categories.json")
	if len(file) != 0 {
		_ = json.Unmarshal(file, &data.Categories)
		for style, artists := range data.Categories {
			for _, artist := range artists {
				data.Artists[artist.Id-1].Category = append(data.Artists[artist.Id-1].Category, style)
			}
		}
		storeArtists()
	} else {
		ApiCategoryFill()
		for style, artists := range data.Categories {
			for _, artist := range artists {
				data.Artists[artist.Id-1].Category = append(data.Artists[artist.Id-1].Category, style)
			}
		}
		storeArtists()
	}
}
