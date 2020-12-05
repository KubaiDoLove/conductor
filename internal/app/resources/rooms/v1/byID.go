package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (rr RoomsResource) RoomByID(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	room, err := rr.repo.RoomByID(context.Background(), roomObjID)
	if err != nil {
		_ = render.Render(w, r, httperrors.ResourceNotFound(err))
		return
	}

	render.JSON(w, r, room)
}
