
package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"amass/internal/infra/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string

	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	HealthCheck     time.Duration
}

func Connect(
	ctx context.Context,
	cfg Config,
) (*pgxpool.Pool, error) {

	logger.Info("initializing PostgreSQL connection...")

	// Escape password safely
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database dsn failed: %w", err)
	}

	// Pool tuning (with sane defaults)
	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	} else {
		poolCfg.MaxConns = 40
	}

	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	} else {
		poolCfg.MinConns = 1
	}

	if cfg.HealthCheck > 0 {
		poolCfg.HealthCheckPeriod = cfg.HealthCheck
	} else {
		poolCfg.HealthCheckPeriod = 30 * time.Second
	}

	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	} else {
		poolCfg.MaxConnLifetime = 1 * time.Hour
	}

	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	} else {
		poolCfg.MaxConnIdleTime = 15 * time.Minute
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool failed: %w", err)
	}

	// Verify connection
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info("database connected successfully")

	return pool, nil
}

