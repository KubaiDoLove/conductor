package v1

import (
	"context"
	"encoding/json"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"github.com/go-chi/render"
	"net/http"
)

func (rr RoomsResource) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	updatedRoom := new(models.Room)
	if err := json.NewDecoder(r.Body).Decode(updatedRoom); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if updatedRoom.MaxLandingPercent > 100 {
		updatedRoom.MaxLandingPercent = 100
	}

	if err := rr.repo.Update(context.Background(), updatedRoom); err != nil {
		_ = render.Render(w, r, httperrors.UnprocessableEntity(err))
		return
	}
}
