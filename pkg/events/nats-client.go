package events

import (
	authHandler "filmserver/pkg/events/handlers/auth"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	connection *nats.Conn
}

func New(natsUrl string) *NatsClient {
	connection, err := GetNatsConnection((natsUrl))
	if err != nil {
		panic(err)
	}

	return &NatsClient{connection: connection}
}

func (c *NatsClient) Auth() *authHandler.Auth {
	return authHandler.New(c.connection)
}
