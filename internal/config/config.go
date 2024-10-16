package config

// PGConfig represents configuration for PostgreSQL.
type PGConfig interface {
	GetDSN() string
}

// GRPCConfig represents configuration for gRPC.
type GRPCConfig interface {
	GetAddress() string
}
