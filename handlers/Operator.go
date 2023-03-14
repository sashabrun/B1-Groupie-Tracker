package handlers

import (
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

func removeFav(id int) {
	for index, favArtist := range data.Fav {
		if favArtist.Name == data.Artists[id].Name {
			data.Fav = RemoveIndex(data.Fav, index)
		}
	}
}
func RemoveIndex(s []Artist, index int) []Artist {
	return append(s[:index], s[index+1:]...)
}
