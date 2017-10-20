package main

import (
	"haste/config"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to Haste!"))
	})

	http.ListenAndServe(":"+strconv.Itoa(config.App().Port), router)
}
