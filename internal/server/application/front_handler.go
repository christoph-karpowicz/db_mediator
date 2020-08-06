package application

import (
	"net/http"
	"text/template"
)

type frontHandler struct {
	app *Application
}

func (h *frontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/build/index.html")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
