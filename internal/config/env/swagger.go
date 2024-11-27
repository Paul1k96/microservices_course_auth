package env

import (
	"net"
	"os"
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/config"
	"github.com/pkg/errors"
)

const (
	swaggerHostEnvName              = "SWAGGER_HOST"
	swaggerPortEnvName              = "SWAGGER_PORT"
	swaggerReadHeaderTimeoutEnvName = "SWAGGER_READ_HEADER_TIMEOUT"
	swaggerShutdownTimeoutEnvName   = "SWAGGER_SHUTDOWN_TIMEOUT"
)

type swaggerConfig struct {
	host              string
	port              string
	readHeaderTimeout time.Duration
	shutdownTimeout   time.Duration
}

// NewSwaggerConfig creates a new swagger config
func NewSwaggerConfig() (config.HTTPConfig, error) {
	host := os.Getenv(swaggerHostEnvName)

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv(swaggerShutdownTimeoutEnvName))
	if err != nil {
		return nil, errors.New("failed to parse swagger shutdown timeout")
	}

	readHeaderTimeout, err := time.ParseDuration(os.Getenv(swaggerReadHeaderTimeoutEnvName))
	if err != nil {
		return nil, errors.New("failed to parse swagger read header timeout")
	}

	return &swaggerConfig{
		host:              host,
		port:              port,
		shutdownTimeout:   shutdownTimeout,
		readHeaderTimeout: readHeaderTimeout,
	}, nil
}

// GetAddress returns the address
func (cfg *swaggerConfig) GetAddress() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// GetReadHeaderTimeout returns the read header timeout
func (cfg *swaggerConfig) GetReadHeaderTimeout() time.Duration {
	return cfg.readHeaderTimeout
}

// GetGracefulShutdownTimeout returns the graceful shutdown timeout
func (cfg *swaggerConfig) GetGracefulShutdownTimeout() time.Duration {
	return cfg.shutdownTimeout
}
