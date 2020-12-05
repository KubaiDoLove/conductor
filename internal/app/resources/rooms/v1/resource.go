package v1

import (
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"github.com/go-chi/chi"
)

type RoomsResource struct {
	repo drivers.RoomsRepository
}

func NewRoomsResource(repo drivers.RoomsRepository) *RoomsResource {
	return &RoomsResource{repo: repo}
}

func (rr RoomsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/new", rr.NewRoom)
	r.Put("/", rr.UpdateRoom)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rr.RoomByID)
		r.Delete("/", rr.DeleteRoom)
	})

	return r
}
