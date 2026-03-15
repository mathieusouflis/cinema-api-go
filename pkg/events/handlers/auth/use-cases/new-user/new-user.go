package newUser

import "github.com/nats-io/nats.go"

type NewUser struct {
	connection *nats.Conn
}

type NewUserPublishInput struct {
	UserId string `json:"user_id"`
}

func New(connection *nats.Conn) *NewUser {
	return &NewUser{connection: connection}
}

func (c *NewUser) Publish(data *NewUserPublishInput) {
	c.connection.Publish("auth.new", []byte(data.UserId))
}
func (c *NewUser) Subscribe(callback nats.MsgHandler) {
	c.connection.Subscribe("auth.new", callback)
}
