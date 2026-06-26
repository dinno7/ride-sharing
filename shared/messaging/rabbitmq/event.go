package rabbitmq

type MessageInfoKind string

const (
	MessageInfoKindEvent   MessageInfoKind = "event"
	MessageInfoKindCommand MessageInfoKind = "command"
)

type (
	MessageInfo[T any] struct {
		ID         string          `json:"id"`
		OwnerID    string          `json:"owner_id"`
		RoutingKey string          `json:"routing_key"`
		Kind       MessageInfoKind `json:"kind"`
		Source     string          `json:"source"` // which service published
		Data       T               `json:"data"`
	}
)
