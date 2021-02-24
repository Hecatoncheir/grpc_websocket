package socket

import (
	"context"
	"fmt"
	"github.com/Hecatoncheir/grpc_websocket/connection"
	"github.com/Hecatoncheir/grpc_websocket/envelope"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	socketGrpc "github.com/Hecatoncheir/grpc_websocket/grpc"

	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"

	"github.com/gorilla/websocket"
)

func TestServer_GRPCRun(t *testing.T) {
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	server := socketGrpc.New()

	go func() {
		grpcIp := "127.0.0.1"
		grpcPort, _, err := GetFreePort()
		if err != nil {
			t.Error(err)
		}

		grpcServerEndpoint := fmt.Sprintf("%v:%v", grpcIp, grpcPort)

		go func() {
			listener, err := net.Listen("tcp", grpcServerEndpoint)
			if err != nil {
				t.Error(err)
			}

			grpcServer := grpc.NewServer()

			socketGrpc.RegisterServiceServer(
				grpcServer,
				server,
			)

			reflection.Register(grpcServer)
			err = grpcServer.Serve(listener)
			if err != nil {
				t.Error(err)
			}

		}()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err = socketGrpc.RegisterServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
		if err != nil {
			t.Error(err)
		}

		testServer := httptest.NewServer(wsproxy.WebsocketProxy(mux))
		defer testServer.Close()

		// Convert http://127.0.0.1 to ws://127.0.0.
		address := "ws" + strings.TrimPrefix(testServer.URL, "http")
		address = address + "/ws"

		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(address, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		defer func() {
			err := ws.Close()
			if err != nil {
				t.Error(err)
			}
		}()

		var connectedClient connection.ConnectionInterface

		for connect := range server.GetOnConnectChannel() {
			if connect == nil {
				t.Fail()
			} else {
				connectedClient = connect
				break
			}
		}

		if connectedClient == nil {
			t.Fail()
			return
		}

		event := envelope.Envelope{
			Message: "Test message from server",
			Details: "",
		}

		err = connectedClient.WriteMessage(&event)
		if err != nil {
			t.Error(err)
		}

		_, message, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if string(message) != "Test message from server" {
			t.Fail()
		}

		//if err := ws.WriteMessage(websocket.TextMessage, []byte("Test message from client")); err != nil {
		//	t.Fatalf("%v", err)
		//}
		//
		//for event := range connectedClient.GetInputChan() {
		//	if event.GetMessage() != "Test message from client" {
		//		t.Fail()
		//	} else {
		//		break
		//	}
		//}

		//err = ws.Close()
		//if err != nil {
		//	t.Error(err)
		//}
		//
		//for event := range connectedClient.GetDisconnectChan() {
		//	if event != true {
		//		t.Fail()
		//	} else {
		//		break
		//	}
		//}

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}

}

func GetFreePort() (int, *net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil, err
	}

	err = l.Close()
	if err != nil {
		return 0, nil, err
	}

	return l.Addr().(*net.TCPAddr).Port, l, nil
}
