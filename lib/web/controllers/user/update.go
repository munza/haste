package userctrl

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte("update user " + userID + "!"))
}
