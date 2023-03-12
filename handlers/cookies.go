package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
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
	return valueToSend
}
