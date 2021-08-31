package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePlanets(t *testing.T) {
	t.Run("Test if it can create planets", func(t *testing.T) {
		tt := []struct {
			name       string
			method     string
			body       string
			want       string
			statusCode int
		}{
			{
				name:       "with empty name",
				method:     http.MethodPost,
				body:       `{"name":"","climate":"arid", "terrain":"stone"}`,
				want:       `Name can't be empty`,
				statusCode: http.StatusUnprocessableEntity,
			},
			{
				name:       "with empty climate",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"", "terrain":"stone"}`,
				want:       `Climate can't be empty`,
				statusCode: http.StatusUnprocessableEntity,
			},
			{
				name:       "with empty terrain",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"arid", "terrain":""}`,
				want:       `Terrain can't be empty`,
				statusCode: http.StatusUnprocessableEntity,
			},
			{
				name:       "with valid data",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"arid", "terrain":"test"}`,
				want:       ``,
				statusCode: http.StatusCreated,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				request, _ := http.NewRequest(tc.method, "/planets", strings.NewReader(tc.body))
				response := httptest.NewRecorder()

				s := &StorageMem{}
				swapi := SwapiMem{}
				h := NewPlanetHandler(s, swapi)

				h.CreatePlanetHandler(response, request)

				got := response.Body.String()
				want := tc.want

				if got != want {
					t.Errorf("Expected msg %q, got %q", want, got)
				}

				gotCode := response.Code
				wantCode := tc.statusCode

				if gotCode != wantCode {
					t.Errorf("Expected status %d, got %d", wantCode, gotCode)
				}
			})
		}

	})

	t.Run("Test if it get films from swapi", func(t *testing.T) {
		tt := []struct {
			name       string
			method     string
			body       string
			want       int
			statusCode int
		}{
			{
				name:       "with 2 films",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"arid", "terrain":"stone"}`,
				want:       2,
				statusCode: http.StatusCreated,
			},
			{
				name:       "with 3 films",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"arid", "terrain":"stone"}`,
				want:       3,
				statusCode: http.StatusCreated,
			},
			{
				name:       "with 0 films",
				method:     http.MethodPost,
				body:       `{"name":"test","climate":"arid", "terrain":"stone"}`,
				want:       0,
				statusCode: http.StatusCreated,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				request, _ := http.NewRequest(tc.method, "/planets", strings.NewReader(tc.body))
				response := httptest.NewRecorder()

				s := &StorageMem{}
				swapi := SwapiMem{tc.want}
				h := NewPlanetHandler(s, swapi)

				h.CreatePlanetHandler(response, request)

				got := s.Planets[0].FilmCounter
				want := tc.want

				if got != want {
					t.Errorf("Expected planet %d, got %d", want, got)
				}

				gotCode := response.Code
				wantCode := tc.statusCode

				if gotCode != wantCode {
					t.Errorf("Expected status %d, got %d", wantCode, gotCode)
				}
			})
		}
	})

	t.Run("Test if it can list planets", func(t *testing.T) {
		tt := []struct {
			name       string
			method     string
			body       string
			want       string
			statusCode int
			counter    int
		}{
			{
				name:       "with 1 films",
				method:     http.MethodGet,
				body:       `{}`,
				want:       `[{"ID":"000000000000000000000000","Name":"SandPlanet","Climate":"Tropical","Terrain":"sand","FilmCounter":0}]`,
				statusCode: http.StatusOK,
				counter:    1,
			},
			{
				name:       "with 2 films",
				method:     http.MethodGet,
				body:       `{}`,
				want:       `[{"ID":"000000000000000000000000","Name":"SandPlanet","Climate":"Tropical","Terrain":"sand","FilmCounter":0},{"ID":"000000000000000000000000","Name":"SandPlanet","Climate":"Tropical","Terrain":"sand","FilmCounter":0}]`,
				statusCode: http.StatusOK,
				counter:    2,
			},
			{
				name:       "with 0 films",
				method:     http.MethodGet,
				body:       `{}`,
				want:       `[]`,
				statusCode: http.StatusOK,
				counter:    0,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				request, _ := http.NewRequest(tc.method, "/planets", strings.NewReader(tc.body))
				response := httptest.NewRecorder()

				s := &StorageMem{make([]Planet, 0)}
				for i := 0; i < tc.counter; i++ {
					s.Planets = append(s.Planets, Planet{Name: "SandPlanet", Climate: "Tropical", Terrain: "sand"})
				}
				swapi := SwapiMem{}

				h := NewPlanetHandler(s, swapi)

				h.GetPlanetsHandler(response, request)

				got := response.Body.String()
				want := tc.want

				if got != want {
					t.Errorf("Expected return %s, got %s", want, got)
				}

				gotCode := response.Code
				wantCode := tc.statusCode

				if gotCode != wantCode {
					t.Errorf("Expected status %d, got %d", wantCode, gotCode)
				}
			})
		}

	})

}

type SwapiMem struct {
	Count int
}

func (s SwapiMem) GetFilmCount(name string) (int, error) {
	return s.Count, nil
}

type StorageMem struct {
	Planets []Planet
}

func (s *StorageMem) Add(p Planet) error {
	s.Planets = append(s.Planets, p)

	return nil
}

func (s StorageMem) Remove(id string) error {
	return nil
}

func (s StorageMem) Find(f []Filter) ([]Planet, error) {
	return s.Planets, nil
}
