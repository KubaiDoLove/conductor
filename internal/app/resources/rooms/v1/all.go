package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/render"
	"net/http"
)

func (rr RoomsResource) All(w http.ResponseWriter, r *http.Request) {
	rooms, err := rr.repo.All(context.Background())
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	render.JSON(w, r, rooms)
}
