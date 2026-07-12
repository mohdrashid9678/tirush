package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Service holds the connection pool
type Service struct {
	Db *pgxpool.Pool
}

// New initializes the database connection
func New(connStr string) (*Service, error) {

	// Parse Config
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	// Pool Configuration
	// MaxConns: Max number of connections in the pool.
	// If all are busy, new queries wait.
	dbConfig.MaxConns = 50
	dbConfig.MinConns = 5
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return &Service{Db: pool}, nil
}

// Close closes the connection pool
func (s *Service) Close() {
	if s.Db != nil {
		s.Db.Close()
		log.Println("Database connection closed")
	}
}
