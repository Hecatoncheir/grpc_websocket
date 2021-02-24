package configuration

type ConfigurationInterface interface {
	GetGrpcSocketServerIp() string
	GetGrpcSocketServerPort() int

	GetSocketServerIp() string
	GetSocketServerPort() int
}
