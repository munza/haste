package userctrl

import (
	auth "haste/lib/auth/repos"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func Show(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
	if user, err := auth.UserRepo().Find(userID); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("show user " + user.Name + "!"))
	}
}
