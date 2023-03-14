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

func DisplayLocationLink(artistId int) string {
	apiKey := "AIzaSyBop5uo9b8uRNxiLd8WbK0ep0yS9ltu7K8"
	address := "Toulouse"
	urlEncodedAddress := url.QueryEscape(address)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", urlEncodedAddress, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status %s", resp.Status)
	}

	var geocodeResponse GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&geocodeResponse)
	if err != nil {
		log.Fatalf("JSON decoding failed: %v", err)
	}

	if resp.Status != "OK" {
		fmt.Println("Geocoding API returned status %s", resp.Status)
	}

	if len(geocodeResponse.Results) == 0 {
		log.Fatalf("No geocode results found")
	}

	lat := geocodeResponse.Results[0].Geometry.Location.Lat
	lng := geocodeResponse.Results[0].Geometry.Location.Lng
	//fmt.Println("Latitude: %f, Longitude: %f", lat, lng)

	return "https://www.google.com/maps/place/" + strconv.FormatFloat(lat, 'f', 9, 64) + "," + strconv.FormatFloat(lng, 'f', 9, 64)
}
