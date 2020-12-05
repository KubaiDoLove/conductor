package drivers

import "errors"

var ErrInvalidConfigStruct = errors.New("invalid config structure")

var (
	ErrEmptyRoom        = errors.New("empty room struct")
	ErrRoomDoesNotExist = errors.New("room does not exist")
)

var (
	ErrEmptyPlace        = errors.New("empty place struct")
	ErrPlaceDoesNotExist = errors.New("place does not exist")
	ErrPlaceTaken        = errors.New("place is already taken")
	ErrInvalidPlace      = errors.New("invalid place")
)
