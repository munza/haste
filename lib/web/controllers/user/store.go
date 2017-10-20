package userctrl

import (
	"net/http"
)

func Store(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user!"))
}
