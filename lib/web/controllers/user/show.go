package userctrl

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Show(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte("show user " + userID + "!"))
}
