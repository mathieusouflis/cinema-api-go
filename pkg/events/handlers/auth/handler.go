package authHandler

import (
	newUser "filmserver/pkg/events/handlers/auth/use-cases/new-user"

	"github.com/nats-io/nats.go"
)

type Auth struct {
	connection *nats.Conn
}

func New(connection *nats.Conn) *Auth {
	return &Auth{connection: connection}
}

func (a *Auth) Create() *newUser.NewUser {
	return newUser.New(a.connection)
}
