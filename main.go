package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	url := os.Getenv("MONGO_URL")
	s, err := NewStorage(url)

	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	swapi := Swapi{}
	h := NewPlanetHandler(s, swapi)

	r.HandleFunc("/planets", h.GetPlanetsHandler).Methods(http.MethodGet)
	r.HandleFunc("/planets", h.CreatePlanetHandler).Methods(http.MethodPost)
	r.HandleFunc("/planets/{id}", h.RemovePlanetHandler).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}

	log.Println("Server listening on :8080...")
	log.Fatal(srv.ListenAndServe())
}
