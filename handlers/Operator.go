package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
