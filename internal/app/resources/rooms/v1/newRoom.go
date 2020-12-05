package v1

import (
	"context"
	"encoding/json"
	httperrors "github.com/KubaiDoLove/conductor/internal/app/errors/http"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"github.com/go-chi/render"
	"net/http"
)

type NewRoomRequest struct {
	Name              string                `json:"name"`
	Size              models.RoomSize       `json:"roomSize"`
	MaxLandingPercent *uint8                `json:"maxLandingPercent"`
	SocialDistance    *uint                 `json:"socialDistance"`
	WorkingHours      models.OpenCloseHours `json:"workingHours"`
}

func (rr RoomsResource) NewRoom(w http.ResponseWriter, r *http.Request) {
	request := new(NewRoomRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	if request.MaxLandingPercent == nil {
		maxPercent := uint8(100)
		request.MaxLandingPercent = &maxPercent
	}

	if request.SocialDistance == nil {
		zeroSocialDistance := uint(0)
		request.SocialDistance = &zeroSocialDistance
	}

	room := models.NewRoom(request.Name, request.Size, *request.MaxLandingPercent, *request.SocialDistance, request.WorkingHours)
	if err := rr.repo.Create(context.Background(), room); err != nil {
		_ = render.Render(w, r, httperrors.UnprocessableEntity(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
