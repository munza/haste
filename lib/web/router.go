package web

import (
	"net/http"

	"haste/lib/web/controllers/user"

	"github.com/go-chi/chi"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", userctrl.Index)
	r.Post("/", userctrl.Store)
	r.Route("/{userID}", func(r chi.Router) {
		r.Get("/", userctrl.Show)
		r.Put("/", userctrl.Update)
		r.Delete("/", userctrl.Delete)
	})

	return r
}
