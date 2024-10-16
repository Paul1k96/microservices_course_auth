package env

import (
	"net"
	"os"

	"github.com/Paul1k96/microservices_course_auth/internal/config"
)

const (
	grpcHost = "GRPC_HOST"
	grpcPort = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig returns a new config.GRPCConfig.
func NewGRPCConfig() config.GRPCConfig {
	var cfg grpcConfig

	cfg.host = os.Getenv(grpcHost)
	cfg.port = os.Getenv(grpcPort)

	return &cfg
}

// GetAddress returns the address of the gRPC server.
func (c *grpcConfig) GetAddress() string {
	return net.JoinHostPort(c.host, c.port)
}
