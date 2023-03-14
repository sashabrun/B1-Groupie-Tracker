package handlers

import (
	"encoding/json"
	"net/http"
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
}
func FillLocation() {
	for _, artist := range data.Artists {
		relationsResp, _ := http.Get(artist.RelationsLink)
		_ = json.NewDecoder(relationsResp.Body).Decode(&artist.Relations)
		defer relationsResp.Body.Close()
	}
}

/* func DisplayLocationLink(artistId int) string {
	return "https://www.google.com/maps/place/" + strconv.FormatFloat(lat, 'f', 9, 64) + "," + strconv.FormatFloat(lng, 'f', 9, 64)
}*/
