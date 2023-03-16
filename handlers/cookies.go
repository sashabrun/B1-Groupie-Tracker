package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DecodeFavCookie(cookie *http.Cookie) {
	cookieData, _ := base64.URLEncoding.DecodeString(cookie.Value)
	_ = json.Unmarshal(cookieData, &data.FavIndexs)
	for index, _ := range data.Artists {
		data.Artists[index].Isliked = false
	}
	for _, index := range data.FavIndexs {
		data.Artists[index].Isliked = true
	}
}

func EncodeFavCookieValue() string {
	FavJSON, _ := json.Marshal(data.FavIndexs)
	//Encoding the json string as a b64 URL before sending it to the cookie to prevent unhandeled characters json uses
	valueToSend := base64.URLEncoding.EncodeToString(FavJSON)
	return valueToSend
}
func CheckFavCookie(r *http.Request) {
	if favCookie, err := r.Cookie("Fav"); err == nil {
		DecodeFavCookie(favCookie)
	} else {
		fmt.Println("No \"Fav\" cookie yet")
	}
}
func UpdateFavCookie(w http.ResponseWriter) {
	favCookie := &http.Cookie{
		Name:  "Fav",
		Value: EncodeFavCookieValue(),
		//The "Fav" cookie has to never expire
		//to save client's stats in his navigator :
		//their favorite artists.
		Expires: time.Date(2024, 12, 01, 00, 00, 00, 00, time.UTC),
		Path:    "/",
	}
	http.SetCookie(w, favCookie)
}
