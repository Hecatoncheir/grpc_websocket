package connected_client

import (
	"github.com/Hecatoncheir/grpc_websocket/connection"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ConnectedClientInterface

	id         string
	connection connection.ConnectionInterface

	httpConnection *websocket.Conn
}

func New(connection connection.ConnectionInterface) *Client {
	id := uuid.New().String()

	client := Client{
		id:         id,
		connection: connection,
	}

	return &client
}

func (client *Client) GetID() string {
	return client.id
}

func (client *Client) GetConnection() connection.ConnectionInterface {
	return client.connection
}

func (client *Client) Disconnect() error {
	client.connection.GetDisconnectChan() <- true
	return nil
}

func (client *Client) GetMessagesChan() <-chan envelope.EnvelopeInterface {
	return client.connection.GetInputChan()
}

func (client *Client) WriteMessage(envelope envelope.EnvelopeInterface) error {
	err := client.connection.WriteMessage(envelope)
	if err != nil {
		return err
	}

	return nil
}
