package domain

import "errors"

var (
	ErrInvalidTripStatus = errors.New("invalid trip status")
	ErrTripNotFound      = errors.New("trip not found")

	ErrFareNotFound = errors.New("fare not found")
)
