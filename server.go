package socket

import (
	"context"
	"fmt"

	"github.com/Hecatoncheir/grpc_websocket/configuration"
	"github.com/Hecatoncheir/grpc_websocket/connected_client"
	socketGrpc "github.com/Hecatoncheir/grpc_websocket/grpc"

	"log"
	"net/http"
	"os"
	"sync"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
)

type Server struct {
	ServerInterface

	config configuration.ConfigurationInterface

	clients   map[string]connected_client.ConnectedClientInterface
	clientsMu sync.Mutex

	onClientAdded chan connected_client.ConnectedClientInterface

	/// TODO: add onClientDisconnected

	logger *log.Logger
}

func New(config configuration.ConfigurationInterface) (*Server, error) {
	logger := log.New(os.Stdout, "SocketServer: ", log.LstdFlags)

	onClientAdded := make(chan connected_client.ConnectedClientInterface)

	engine := Server{
		config:        config,
		logger:        logger,
		onClientAdded: onClientAdded,
	}

	return &engine, nil
}

func (server *Server) Run() error {
	return server.runWithGRPC()
}

func (server *Server) runWithGRPC() error {
	var err error

	grpcIp := server.config.GetGrpcSocketServerIp()
	grpcPort := server.config.GetGrpcSocketServerPort()
	grpcServerEndpoint := fmt.Sprintf("%v:%v", grpcIp, grpcPort)

	go func() {
		grpcServer, err := socketGrpc.RunGRPCServer(grpcServerEndpoint)
		if err != nil {
			server.logger.Println(err)
		}

		go func() {
			connectedClientChan := grpcServer.GetOnConnectChannel()
			for connection := range connectedClientChan {
				client := connected_client.New(connection)
				err := server.AddClient(client)
				if err != nil {
					server.logger.Println(err)
				}

				logMessage := fmt.Sprintf("New client: %v added", client.GetID())
				server.logger.Println(logMessage)
			}
		}()

		go func() {
			disconnectedClientChan := grpcServer.GetOnDisconnectChannel()
			for connection := range disconnectedClientChan {

				var client connected_client.ConnectedClientInterface
				for key := range server.clients {
					connectedClient := server.clients[key]
					if connectedClient.GetConnection() == connection {
						client = connectedClient
						break
					}
				}

				if client == nil {
					logMessage := fmt.Sprintf("Can not find connected client for disconnect")
					server.logger.Println(logMessage)
				} else {

					err := server.RemoveClient(client)
					if err != nil {
						server.logger.Println(err)
					}

					logMessage := fmt.Sprintf("Client: %v removed", client.GetID())
					server.logger.Println(logMessage)
				}
			}
		}()
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = socketGrpc.RegisterServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		server.logger.Println(err)
	}

	socketServerIp := server.config.GetSocketServerIp()
	socketServerPort := server.config.GetSocketServerPort()
	socketServerEndpoint := fmt.Sprintf("%v:%v", socketServerIp, socketServerPort)

	server.logger.Println("Server is ready to handle requests at", socketServerEndpoint)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	err = http.ListenAndServe(socketServerEndpoint, wsproxy.WebsocketProxy(mux))
	if err != nil {
		server.logger.Println(err)
	}

	return nil
}

func (server *Server) AddClient(client connected_client.ConnectedClientInterface) error {

	server.clientsMu.Lock()
	server.clients[client.GetID()] = client
	server.clientsMu.Unlock()

	server.onClientAdded <- client

	return nil
}

func (server *Server) OnClientAddedChannel() chan connected_client.ConnectedClientInterface {
	return server.onClientAdded
}

func (server *Server) RemoveClient(client connected_client.ConnectedClientInterface) error {

	server.clientsMu.Lock()
	delete(server.clients, client.GetID())
	server.clientsMu.Unlock()

	server.OnClientRemovedChannel() <- client

	return nil
}

func (server *Server) OnClientRemovedChannel() chan connected_client.ConnectedClientInterface {
	return server.onClientAdded
}
