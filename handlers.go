package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Repository interface {
	Add(p Planet) error
	Remove(id string) error
	Find(f []Filter) ([]Planet, error)
}

type SwapiSDK interface {
	GetFilmCount(name string) (int, error)
}

type PlanetHandler struct {
	Repo  Repository
	Swapi SwapiSDK
}

func NewPlanetHandler(r Repository, s SwapiSDK) *PlanetHandler {
	return &PlanetHandler{r, s}
}

func (h PlanetHandler) GetPlanetsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var filters []Filter

	if name := query.Get("name"); name != "" {
		filters = append(filters, Filter{Key: "name", Value: name})
	}

	if id := query.Get("id"); id != "" {
		filters = append(filters, Filter{Key: "id", Value: id})
	}

	planets, err := h.Repo.Find(filters)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
	}

	planetsJson, err := json.Marshal(planets)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(planetsJson)
}

func (h PlanetHandler) CreatePlanetHandler(w http.ResponseWriter, r *http.Request) {
	var p Planet
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p.Name == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Name can't be empty"))
		return
	}

	if p.Climate == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Climate can't be empty"))
		return
	}

	if p.Terrain == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Terrain can't be empty"))
		return
	}

	filmCount, err := h.Swapi.GetFilmCount(p.Name)

	if err != nil {
		log.Println(err)
	}

	p.FilmCounter = filmCount
	err = h.Repo.Add(p)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h PlanetHandler) RemovePlanetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, ok := params["id"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Planet id is required"))
		return
	}

	err := h.Repo.Remove(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
