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

	return r
}
