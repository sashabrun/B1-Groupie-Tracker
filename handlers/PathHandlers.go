package handlers

import (
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_ = tpl.ExecuteTemplate(w, "home.html", nil)
}
