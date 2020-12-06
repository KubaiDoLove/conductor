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

	r.Post("/forRoom/{roomID}", pr.NewPlace)

	r.Route("/{placeID}", func(r chi.Router) {
		r.Delete("/", pr.DeletePlace)
		r.Get("/suggestions", pr.Suggestions)
	})

	r.Route("/booking", func(r chi.Router) {
		r.Post("/new", pr.NewBooking)
		r.Delete("/{bookingID}", pr.CancelBooking)
	})

	return r
}
