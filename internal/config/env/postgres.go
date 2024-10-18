package env

import (
	"fmt"
	"os"

	"github.com/Paul1k96/microservices_course_auth/internal/config"
)

// Environment variables.
const (
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD" // nolint: gosec
	postgresDBName   = "POSTGRES_DB"
)

type pgConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

// NewPGConfig returns a new config.PGConfig.
func NewPGConfig() config.PGConfig {
	var cfg pgConfig

	cfg.host = os.Getenv(postgresHost)
	cfg.port = os.Getenv(postgresPort)
	cfg.user = os.Getenv(postgresUser)
	cfg.password = os.Getenv(postgresPassword)
	cfg.dbName = os.Getenv(postgresDBName)

	return &cfg
}

// GetDSN returns the data source name.
func (c *pgConfig) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.user, c.password, c.dbName)

	return dsn
}
