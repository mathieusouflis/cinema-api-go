package events

import "github.com/nats-io/nats.go"

var natsConnection map[string]*nats.Conn

func GetNatsConnection(natsUrl string) (*nats.Conn, error) {
	if natsConnection == nil {
		natsConnection = make(map[string]*nats.Conn)
	}

	if connection := natsConnection[natsUrl]; connection != nil {
		return connection, nil
	}

	connection, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	natsConnection[natsUrl] = connection
	return connection, nil
}
