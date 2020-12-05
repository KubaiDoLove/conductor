package v1

import (
	"context"
	"encoding/json"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type NewBookingRequest struct {
	PlaceID   primitive.ObjectID `json:"placeID"`
	StartTime time.Time          `json:"startTime"`
	EndTime   time.Time          `json:"endTime"`
}

func (pr PlacesResource) NewBooking(w http.ResponseWriter, r *http.Request) {
	request := new(NewBookingRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	booking := models.NewBooking(request.StartTime, request.EndTime)
	if err := pr.repo.AddBooking(context.Background(), request.PlaceID, booking); err != nil {
		_ = render.Render(w, r, httperrors.Conflict(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
