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

func (pr PlacesResource) NewBooking(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeID")
	placeObjID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	booking := new(models.Booking)
	if err := json.NewDecoder(r.Body).Decode(booking); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if err := pr.repo.AddBooking(context.Background(), placeObjID, booking); err != nil {
		_ = render.Render(w, r, httperrors.Conflict(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
