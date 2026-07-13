package domain

import (
	"encoding/json"
	"slices"
)

type TripStatus struct {
	value string
}

var (
	TripStatusPending  = &TripStatus{"pending"}
	TripStatusComplete = &TripStatus{"complete"}
	TripStatusAccepted = &TripStatus{"accepted"}
)

func (ts *TripStatus) String() string {
	return ts.value
}

// INFO: For showing it in json as string not a struct
func (ts *TripStatus) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}
	return json.Marshal(ts.value)
}

func (ts *TripStatus) IsValid() bool {
	return slices.Contains([]*TripStatus{
		TripStatusPending,
		TripStatusAccepted,
		TripStatusComplete,
	}, ts)
}
