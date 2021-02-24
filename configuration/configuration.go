package configuration

import (
	"github.com/caarlos0/env"
)

// Config - configuration values from environment variables.
type Config struct {
	ConfigurationInterface

	GRPCSocketIp   string `env:"GRPCSocketIp" envDefault:"127.0.0.1"`
	GRPCSocketPort int    `env:"GRPCSocketPort" envDefault:"9001"`

	SocketIp   string `env:"SocketIp" envDefault:"127.0.0.1"`
	SocketPort int    `env:"SocketPort" envDefault:"8001"`
}

func New() (*Config, error) {
	configuration := Config{}

	err := env.Parse(&configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (config *Config) GetSocketServerIp() string {
	return config.SocketIp
}

func (config *Config) GetSocketServerPort() int {
	return config.SocketPort
}

func (config *Config) GetGrpcSocketServerIp() string {
	return config.GRPCSocketIp
}

func (config *Config) GetGrpcSocketServerPort() int {
	return config.GRPCSocketPort
}
