package connected_client

import (
	"github.com/Hecatoncheir/grpc_websocket/connection"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
)

type ConnectedClientInterface interface {
	GetID() string

	GetConnection() connection.ConnectionInterface

	Disconnect() error

	GetMessagesChan() <-chan envelope.EnvelopeInterface
	WriteMessage(envelope envelope.EnvelopeInterface) error
}
