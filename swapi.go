package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type swapiResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []planetSwapi
}

type planetSwapi struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	URL            string    `json:"url"`
}

const URL = "https://swapi.dev/api/planets"

type Swapi struct{}

func (sw Swapi) GetFilmCount(planetName string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("%s?search=%s", URL, url.QueryEscape(planetName)))

	if err != nil {
		return 0, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	var sr swapiResponse

	err = json.Unmarshal(responseData, &sr)

	if err != nil {
		return 0, err
	}

	if len(sr.Results) == 0 {
		return 0, nil
	}

	return len(sr.Results[0].Films), nil
}
