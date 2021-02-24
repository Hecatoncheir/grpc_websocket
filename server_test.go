package socket

import (
	"github.com/Hecatoncheir/grpc_websocket/configuration"
	"testing"
)

func TestNew(t *testing.T) {
	config, err := configuration.New()
	if err != nil {
		t.Error(err)
	}

	server, err := New(config)
	if err != nil {
		t.Error(err)
	}

	if server.config.GetSocketServerIp() != "127.0.0.1" {
		t.Fail()
	}

	if server.config.GetSocketServerPort() != 8001 {
		t.Fail()
	}
}
