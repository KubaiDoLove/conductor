package drivers

import "errors"

var ErrInvalidConfigStruct = errors.New("invalid config structure")

var ErrEmptyRoom = errors.New("empty room struct")
var ErrRoomDoesNotExist = errors.New("room does not exist")
