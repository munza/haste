package main

import (
	"fmt"
	"haste/config"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	cfgdb := config.Database()
	fmt.Println(cfgdb.Port)

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to haste!"))
	})

	http.ListenAndServe(":3000", router)
}
