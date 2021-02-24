package connection

import "github.com/Hecatoncheir/grpc_websocket/envelope"

type ConnectionInterface interface {
	GetInputChan() chan envelope.EnvelopeInterface
	GetOutputChan() chan envelope.EnvelopeInterface
	GetDisconnectChan() chan bool
	Disconnect() error
	WriteMessage(envelope envelope.EnvelopeInterface) error
}
