package socket

import (
	"github.com/Hecatoncheir/grpc_websocket/connected_client"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
)

type ServerInterface interface {
	Run() error

	OnClientAddedChannel() chan connected_client.ConnectedClientInterface
	AddClient(client connected_client.ConnectedClientInterface) error

	OnClientRemovedChannel() chan connected_client.ConnectedClientInterface
	RemoveClient(client connected_client.ConnectedClientInterface) error

	GetClients() []connected_client.ConnectedClientInterface
	WriteTo([]connected_client.ConnectedClientInterface, envelope.EnvelopeInterface) error
}
