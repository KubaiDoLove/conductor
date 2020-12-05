package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (pr PlacesResource) DeletePlace(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeID")
	placeObjID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if err := pr.repo.Delete(context.Background(), placeObjID); err != nil {
		_ = render.Render(w, r, httperrors.ResourceNotFound(err))
		return
	}
}
