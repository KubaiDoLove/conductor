package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (rr RoomsResource) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if err := rr.repo.Delete(context.Background(), roomObjID); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}
}
