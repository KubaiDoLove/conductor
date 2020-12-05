package v1

import (
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"github.com/go-chi/chi"
)

type PlacesResource struct {
	repo drivers.PlacesRepository
}

func NewPlacesResource(repo drivers.PlacesRepository) *PlacesResource {
	return &PlacesResource{repo: repo}
}

func (pr PlacesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/forRoom/{roomID}", func(r chi.Router) {
		r.Post("/", pr.NewPlace)
		r.Delete("/{placeID}", pr.DeletePlace)
	})

	return r
}
