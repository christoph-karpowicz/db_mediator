package handlers

import (
	"fmt"
	"net/http"
)

func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
}
