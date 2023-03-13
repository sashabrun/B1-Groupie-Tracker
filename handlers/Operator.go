package handlers

import "strings"

func ArtistNameContainsInput(name string) bool {
	return strings.Contains(strings.ToUpper(name), strings.ToUpper(data.Input.text))
}
