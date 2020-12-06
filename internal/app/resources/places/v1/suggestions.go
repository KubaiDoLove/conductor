package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (pr PlacesResource) Suggestions(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeID")
	placeObjID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	hours, err := pr.repo.Suggestions(context.Background(), placeObjID)
	if err != nil {
		_ = render.Render(w, r, httperrors.ResourceNotFound(err))
		return
	}

	render.JSON(w, r, hours)
}
