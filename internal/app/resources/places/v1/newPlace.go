package v1

import (
	"context"
	"encoding/json"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type NewPlaceRequest struct {
	X            int                   `json:"x"`
	Y            int                   `json:"y"`
	WorkingHours models.OpenCloseHours `json:"workingHours"`
}

func (pr PlacesResource) NewPlace(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	request := new(NewPlaceRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	place := models.NewPlace(request.X, request.Y, request.WorkingHours)
	if err := pr.repo.Create(context.Background(), roomObjID, place); err != nil {
		_ = render.Render(w, r, httperrors.UnprocessableEntity(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
