package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ArtistNameContainsInput(id int) bool {
	id--
	if data.Input.nbMembers != "" {
		strNbMembers := strconv.Itoa(len(data.Artists[id].Members))
		if strings.Contains(strNbMembers, data.Input.nbMembers) {
			return true
		}
	} else if data.Input.creaDate != "" {
		CreaDate, _ := strconv.Atoi(data.Input.creaDate)
		if data.Artists[id].CreationDate >= CreaDate {
			return true
		}
	} else if strings.Contains(strings.ToUpper(data.Artists[id].Name), strings.ToUpper(data.Input.text)) || strings.Contains(strings.ToUpper(data.Artists[id].FirstAlbum), strings.ToUpper(data.Input.text)) || strings.Contains(strings.ToUpper(strconv.Itoa(data.Artists[id].CreationDate)), strings.ToUpper(data.Input.text)) {
		return true
	}
	return false
}

func removeFav(idToDelete int) {
	for i, favIndex := range data.FavIndexs {
		if favIndex == idToDelete {
			RemoveIndex(data.FavIndexs, i)
		}
	}
}
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
func GetArtists() {
	file, _ := os.ReadFile("data/artists.json")
	if len(file) != 0 {
		_ = json.Unmarshal(file, &data.Artists)
	} else {
		apiRes, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
		if err == nil {
			_ = json.NewDecoder(apiRes.Body).Decode(&data.Artists)
		}
		defer apiRes.Body.Close()
		storeArtists()
	}
}
func FillData() {
	GetArtists()
	FillLocation()
	GetLikes()
}

func FillLocation() {
	for id, artist := range data.Artists {
		relationsResp, _ := http.Get(artist.RelationsLink)
		_ = json.NewDecoder(relationsResp.Body).Decode(&data.Artists[id].Relations)
		defer relationsResp.Body.Close()
	}
}

type GeocodeResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func DisplayLocationLink(cityName string) string {
	return "https://www.google.com/maps/search/" + cityName
}
func GetLikes() {
	file, _ := os.ReadFile("data/Likes.json")
	if len(file) != 0 {
		_ = json.Unmarshal(file, &data.Likes)
	} else {
		data.Likes = make([]int, len(data.Artists))
		SaveLikes()
	}
}

func SaveLikes() {
	LikesJSON, _ := json.Marshal(data.Likes)
	if err := os.WriteFile("data/Likes.json", LikesJSON, 0777); err != nil {
		fmt.Println(err)
	}
}
func GetArtistLikes(id int) int {
	return data.Likes[id-1]
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
	} else {
		ApiCategoryFill()
	}
	storeArtists()
}
