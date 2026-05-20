package rabbitmq

const (
	ExchangeMain = "main.topic"
	ExchangeDLX  = "dlx.topic"
)

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
}

// All exist exchanges for this project
var exchanges = []*Exchange{
	{
		Name:       ExchangeMain,
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
	},
}
