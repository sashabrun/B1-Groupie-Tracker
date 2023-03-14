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
	//id, _ := strconv.Atoi(ids)
	if strings.Contains(strings.ToUpper(data.Artists[id].Name), strings.ToUpper(data.Input.text)) == true || strings.Contains(strings.ToUpper(data.Artists[id].FirstAlbum), strings.ToUpper(data.Input.text)) == true || strings.Contains(strings.ToUpper(strconv.Itoa(data.Artists[id].CreationDate)), strings.ToUpper(data.Input.text)) == true {
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

func FillData() {
	apiRes, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err == nil {
		_ = json.NewDecoder(apiRes.Body).Decode(&data.Artists)
	}
	defer apiRes.Body.Close()
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
		data.Likes = make([]int, 52)
		SaveLikes()
	}
}

func SaveLikes() {
	LikesJSON, _ := json.Marshal(data.Likes)
	if err := os.WriteFile("data/Likes.json", LikesJSON, 0777); err != nil {
		fmt.Println(err)
	}
}
