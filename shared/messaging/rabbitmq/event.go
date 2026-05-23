package rabbitmq

import (
	"time"
)

type Event[T any] struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Type      string    `json:"type"`
	Source    string    `json:"source"` // which service published
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}
