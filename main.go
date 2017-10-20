package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// fmt.Println(config.Get("Database.Missing"))

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to haste!"))
	})

	http.ListenAndServe(":3000", router)
}
