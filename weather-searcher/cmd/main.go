package main

import (
	"encoding/json"
	"net/http"
)

type Weather struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		weather := Weather{City: "any_city", TempC: 100.0, TempF: 100.0, TempK: 100.0}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(weather)
	})
	http.ListenAndServe(":8081", nil)
}
