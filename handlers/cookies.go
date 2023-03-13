package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DecodeFavCookie(cookie *http.Cookie) []Artist {
	cookieData, _ := base64.URLEncoding.DecodeString(cookie.Value)
	favs := make([]Artist, 0)
	_ = json.Unmarshal(cookieData, &favs)
	return favs
}

func EncodeFavCookieValue(favs []Artist) string {
	FavJSON, _ := json.Marshal(favs)
	//Encoding the json string as a b64 URL before sending it to the cookie to prevent unhandeled characters json uses
	valueToSend := base64.URLEncoding.EncodeToString(FavJSON)
	fmt.Println(valueToSend)
	return valueToSend
}

func UpdateFavCookie(w http.ResponseWriter) {
	favCookie := &http.Cookie{
		Name:  "Fav",
		Value: EncodeFavCookieValue(data.Artists),
		//The "Fav" cookie has to never expire
		//to save client's stats in his navigator :
		//their favorite artists.
		Expires: time.Date(2037, 12, 01, 00, 00, 00, 00, time.UTC),
	}
	http.SetCookie(w, favCookie)
}
