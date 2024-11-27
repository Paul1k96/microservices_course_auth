package env

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/config"
)

const (
	httpHostEnvName              = "HTTP_HOST"
	httpPortEnvName              = "HTTP_PORT"
	httpShutdownTimeoutEnvName   = "HTTP_SHUTDOWN_TIMEOUT"
	httpReadHeaderTimeoutEnvName = "HTTP_READ_HEADER_TIMEOUT"
)

type httpConfig struct {
	host              string
	port              string
	shutdownTimeout   time.Duration
	readHeaderTimeout time.Duration
}

// NewHTTPConfig returns new HTTP config.
func NewHTTPConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv(httpShutdownTimeoutEnvName))
	if err != nil {
		return nil, fmt.Errorf("failed to parse http shutdown timeout: %w", err)
	}

	readHeaderTimeout, err := time.ParseDuration(os.Getenv(httpReadHeaderTimeoutEnvName))
	if err != nil {
		return nil, fmt.Errorf("failed to parse http read header timeout: %w", err)
	}

	return &httpConfig{
		host:              host,
		port:              port,
		shutdownTimeout:   shutdownTimeout,
		readHeaderTimeout: readHeaderTimeout,
	}, nil
}

// GetAddress returns address.
func (cfg *httpConfig) GetAddress() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// GetReadHeaderTimeout returns read header timeout.
func (cfg *httpConfig) GetReadHeaderTimeout() time.Duration {
	return cfg.readHeaderTimeout
}

// GetGracefulShutdownTimeout returns graceful shutdown timeout.
func (cfg *httpConfig) GetGracefulShutdownTimeout() time.Duration {
	return cfg.shutdownTimeout
}
