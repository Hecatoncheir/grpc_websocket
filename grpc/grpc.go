package grpc

import (
	"fmt"
	"github.com/Hecatoncheir/grpc_websocket/connection"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func RunGRPCServer(grpcServerEndpoint string) (GRPCServerInterface, error) {

	listener, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	logger := log.New(os.Stdout, "GrpcSocketServer: ", log.LstdFlags)
	logger.Println("Server is ready to handle requests at", grpcServerEndpoint)

	server := New()
	RegisterServiceServer(
		grpcServer,
		server,
	)

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		return nil, err
	}

	return server, nil
}

type Server struct {
	GRPCServerInterface
	logger *log.Logger

	onConnectClientChannel    chan connection.ConnectionInterface
	onDisconnectClientChannel chan connection.ConnectionInterface
}

func New() *Server {
	logger := log.New(os.Stdout, "SocketGRPCServer: ", log.LstdFlags)

	onConnectClientChannel := make(chan connection.ConnectionInterface)
	onDisconnectClientChannel := make(chan connection.ConnectionInterface)

	server := Server{
		logger:                    logger,
		onConnectClientChannel:    onConnectClientChannel,
		onDisconnectClientChannel: onDisconnectClientChannel,
	}

	return &server
}

func (server *Server) GetOnConnectChannel() chan connection.ConnectionInterface {
	return server.onConnectClientChannel
}

func (server *Server) GetOnDisconnectChannel() chan connection.ConnectionInterface {
	return server.onDisconnectClientChannel
}

func (server *Server) Stream(stream Service_StreamServer) error {
	server.logger.Println("Client connected")

	inputChannel := make(chan envelope.EnvelopeInterface)
	outputChannel := make(chan envelope.EnvelopeInterface)
	disconnectChannel := make(chan bool)
	closeChannel := make(chan bool)

	connect := connection.NewFromGRPC(inputChannel, outputChannel, disconnectChannel)

	go func(connect connection.ConnectionInterface) {
		server.onConnectClientChannel <- connect
	}(connect)

	go func(stream Service_StreamServer) {
		for {

			event, err := stream.Recv()
			if err != nil {
				disconnectChannel <- true
				break
			}

			logMessage := fmt.Sprintf("Receive message: \"%v\" \n", event.Message)
			server.logger.Println(logMessage)

			envelop := envelope.Envelope{
				Message: event.GetMessage(),
				Details: event.GetDetails(),
			}

			inputChannel <- &envelop
		}

		disconnectChannel <- true
	}(stream)

	go func(stream Service_StreamServer) {
		for envelop := range outputChannel {
			logMessage := fmt.Sprintf("Send message: \"%v\" \n", envelop.GetMessage())
			server.logger.Println(logMessage)

			message := Response{
				Message: envelop.GetMessage(),
				Details: envelop.GetDetails(),
			}

			err := stream.Send(&message)
			if err != nil {
				println("!!!")
				server.logger.Println(err.Error())
				continue
			}
		}
	}(stream)

	go func(connect connection.ConnectionInterface) {
		for event := range disconnectChannel {
			server.onDisconnectClientChannel <- connect
			closeChannel <- event
			close(closeChannel)
			return
		}
	}(connect)

	<-closeChannel

	server.logger.Println("Client disconnected")

	return nil
}
