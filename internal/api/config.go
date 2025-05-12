package api

import (
	"context"
	"sync/atomic"

	"github.com/GitIBB/pursuit/internal/database"
)

type APIConfig struct { // struct to hold configuration for the API
	fileserverHits atomic.Int32      // counter for file server hits
	db             *database.Queries // database connection
	platform       string            // platform name
	jwtSecret      string            // JWT secret for signing tokens
}

func NewAPIConfig(db *database.Queries, platform, jwtSecret string) *APIConfig {
	return &APIConfig{
		fileserverHits: atomic.Int32{},
		db:             db,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}
}

// Controlled Access Methods
func (cfg *APIConfig) IncrementFileserverHits() {
	cfg.fileserverHits.Add(1)
}

func (cfg *APIConfig) ResetFileserverHits() {
	cfg.fileserverHits.Store(0)
}

func (cfg *APIConfig) ResetDatabase(ctx context.Context) error {
	return cfg.db.Reset(ctx)
}

func (cfg *APIConfig) GetPlatform() string {
	return cfg.platform
}

func (cfg *APIConfig) GetJWTSecret() string {
	return cfg.jwtSecret
}
