package main

import (
	"net/http"

	"github.com/ItaloG/go-weather-searcher/internal/handlers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", handlers.GetCepHandler)

	http.ListenAndServe(":8080", r)
}
