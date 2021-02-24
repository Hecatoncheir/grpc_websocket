package grpc

import (
	"context"
	"github.com/Hecatoncheir/grpc_websocket/connection"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var listener *bufconn.Listener
var server *Server

func init() {
	listener = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	server = New()

	RegisterServiceServer(grpcServer, server)

	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}

func TestServer_Stream_CanHandleConnection(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	clientConnection := NewServiceClient(conn)

	_, err = clientConnection.Stream(ctx)
	if err != nil {
		t.Error(err)
	}

	for connect := range server.GetOnConnectChannel() {
		if connect == nil {
			t.Fail()
		} else {
			break
		}
	}
}

func TestServer_Stream_CanHandleDisconnection(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	clientConnection := NewServiceClient(conn)

	_, err = clientConnection.Stream(ctx)
	if err != nil {
		t.Error(err)
	}

	for connect := range server.GetOnConnectChannel() {
		if connect == nil {
			t.Fail()
		} else {
			err := connect.Disconnect()
			if err != nil {
				t.Error(err)
			}
			break
		}
	}

	for event := range server.GetOnDisconnectChannel() {
		if event == nil {
			t.Fail()
		} else {
			break
		}
	}
}

func TestServer_Stream_CanHandleMessages(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	serviceClient := NewServiceClient(conn)

	client, err := serviceClient.Stream(ctx)
	if err != nil {
		t.Error(err)
	}

	request := Request{
		Message: "Test message",
	}

	err = client.Send(&request)
	if err != nil {
		t.Error(err)
	}

	var clientConnection connection.ConnectionInterface

	for connect := range server.GetOnConnectChannel() {
		if connect == nil {
			t.Fail()
		} else {
			clientConnection = connect
			break
		}
	}

	if clientConnection == nil {
		t.Fail()
		return
	}

	var event envelope.EnvelopeInterface

	for clientEvent := range clientConnection.GetInputChan() {
		if clientEvent == nil {
			t.Fail()
		} else {
			event = clientEvent
			break
		}
	}

	if event == nil {
		t.Fail()
		return
	}

	if event.GetMessage() != "Test message" {
		t.Fail()
	}
}

func TestServer_Stream_CanSendMessages(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	serviceClient := NewServiceClient(conn)

	client, err := serviceClient.Stream(ctx)
	if err != nil {
		t.Error(err)
	}

	go func() {

		var clientConnection connection.ConnectionInterface

		for connect := range server.GetOnConnectChannel() {
			if connect == nil {
				t.Fail()
			} else {
				clientConnection = connect
				break
			}
		}

		if clientConnection == nil {
			t.Fail()
			return
		}

		event := envelope.Envelope{
			Message: "Test message",
		}

		err := clientConnection.WriteMessage(&event)
		if err != nil {
			t.Error(err)
		}
	}()

	for {

		req, err := client.Recv()
		if err != nil {
			t.Error(err)
		}

		if req.GetMessage() != "Test message" {
			t.Fail()
		} else {
			break
		}
	}
}

func TestServer_Stream_CanHandleDisconnect(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	serviceClient := NewServiceClient(conn)

	_, err = serviceClient.Stream(ctx)
	if err != nil {
		t.Error(err)
	}

	var clientConnection connection.ConnectionInterface

	for connect := range server.GetOnConnectChannel() {
		if connect == nil {
			t.Fail()
		} else {
			clientConnection = connect
			break
		}
	}

	if clientConnection == nil {
		t.Fail()
		return
	}

	_ = conn.Close()

	isDisconnected := false
	for event := range clientConnection.GetDisconnectChan() {
		if event != true {
			t.Fail()
		} else {
			isDisconnected = true
			break
		}
	}

	if isDisconnected != true {
		t.Fail()
	}
}
