package application

import (
	"net/http"
	"text/template"
)

type frontHandler struct {
	app *Application
}

func (h *frontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("front/build/index.html")
	err := t.Execute(w, "test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
