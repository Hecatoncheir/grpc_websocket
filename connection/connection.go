package connection

import "github.com/Hecatoncheir/grpc_websocket/envelope"

type Connection struct {
	ConnectionInterface

	inputChannel   chan envelope.EnvelopeInterface
	outputChannel  chan envelope.EnvelopeInterface
	disconnectChan chan bool
}

func NewFromGRPC(
	inputChannel chan envelope.EnvelopeInterface,
	outputChannel chan envelope.EnvelopeInterface,
	disconnectChan chan bool,
) *Connection {

	connection := Connection{
		inputChannel:   inputChannel,
		outputChannel:  outputChannel,
		disconnectChan: disconnectChan,
	}

	return &connection
}

func (connection *Connection) GetInputChan() chan envelope.EnvelopeInterface {
	return connection.inputChannel
}

func (connection *Connection) GetOutputChan() chan envelope.EnvelopeInterface {
	return connection.outputChannel
}

func (connection *Connection) GetDisconnectChan() chan bool {
	return connection.disconnectChan
}

func (connection *Connection) Disconnect() error {
	connection.disconnectChan <- true
	return nil
}

func (connection *Connection) WriteMessage(envelope envelope.EnvelopeInterface) error {
	connection.outputChannel <- envelope
	return nil
}
