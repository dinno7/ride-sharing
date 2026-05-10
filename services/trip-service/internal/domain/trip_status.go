package domain

import "encoding/json"

type tripStatus struct {
	value string
}

var (
	TripStatusPending  = &tripStatus{"pending"}
	TripStatusComplete = &tripStatus{"complete"}
)

func (ts *tripStatus) String() string {
	return ts.value
}

func (ts *tripStatus) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}
	return json.Marshal(ts.value)
}

func (ts *tripStatus) IsValid() bool {
	switch ts {
	case TripStatusPending:
		return true
	default:
		return false
	}
}
