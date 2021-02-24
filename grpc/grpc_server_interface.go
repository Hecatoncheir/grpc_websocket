package grpc

import (
	"github.com/Hecatoncheir/grpc_websocket/connection"
)

type GRPCServerInterface interface {
	ServiceServer
	GetOnConnectChannel() chan connection.ConnectionInterface
	GetOnDisconnectChannel() chan connection.ConnectionInterface
}
