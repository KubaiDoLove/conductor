package v1

import (
	"context"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (pr PlacesResource) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := chi.URLParam(r, "bookingID")
	bookingObjID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if err := pr.repo.CancelBooking(context.Background(), bookingObjID); err != nil {
		_ = render.Render(w, r, httperrors.ResourceNotFound(err))
		return
	}
}
