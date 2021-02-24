package configuration

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	configuration, err := New()
	if err != nil {
		t.Error(err)
	}

	if configuration == nil {
		t.Fail()
	}

	if configuration.SocketIp != "127.0.0.1" {
		t.Fail()
	}

	if configuration.SocketPort != 8001 {
		t.Fail()
	}
}

func TestConfig_GetSocketIp(t *testing.T) {
	var err error
	var config *Config

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetSocketServerIp() != "127.0.0.1" {
		t.Fail()
	}

	err = os.Setenv("SocketIp", "localhost")
	if err != nil {
		t.Fail()
	}

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetSocketServerIp() != "localhost" {
		t.Fail()
	}
}

func TestConfig_GetSocketPort(t *testing.T) {
	var err error
	var config *Config

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetSocketServerPort() != 8001 {
		t.Fail()
	}

	err = os.Setenv("SocketPort", "8002")
	if err != nil {
		t.Fail()
	}

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetSocketServerPort() != 8002 {
		t.Fail()
	}
}

func TestConfig_GetGRPCSocketIp(t *testing.T) {
	var err error
	var config *Config

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetGrpcSocketServerIp() != "127.0.0.1" {
		t.Fail()
	}

	err = os.Setenv("GRPCSocketIp", "localhost")
	if err != nil {
		t.Fail()
	}

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetGrpcSocketServerIp() != "localhost" {
		t.Fail()
	}
}

func TestConfig_GetGRPCSocketPort(t *testing.T) {
	var err error
	var config *Config

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetGrpcSocketServerPort() != 9001 {
		t.Fail()
	}

	err = os.Setenv("GRPCSocketPort", "9002")
	if err != nil {
		t.Fail()
	}

	config, err = New()
	if err != nil {
		t.Error(err)
	}

	if config.GetGrpcSocketServerPort() != 9002 {
		t.Fail()
	}
}
