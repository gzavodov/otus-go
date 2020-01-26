package sql

import (
	"context"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/repository"
	"github.com/jackc/pgx/v4"
)

//BaseRepository represents Base PostgreSQL repository
type BaseRepository struct {
	dataSourceName string
	conn           *pgx.Conn
	ctx            context.Context
}

//Connect try to connect to PostgreSQL server
func (r *BaseRepository) Connect() error {
	if r.conn != nil {
		return nil
	}

	if r.dataSourceName == "" {
		return repository.NewError(
			repository.ErrorInvalidConfiguration,
			"empty DSN (data source name)",
		)
	}

	config, err := pgx.ParseConfig(r.dataSourceName)
	if err != nil {
		return repository.WrapErrorf(
			repository.ErrorInvalidConfiguration,
			err,
			"failed to parse DSN (data source name)",
		)
	}

	r.conn, err = pgx.ConnectConfig(r.ctx, config)
	if err != nil {
		return repository.WrapErrorf(
			repository.ErrorInvalidConfiguration,
			err,
			"failed to connect to PostgreSQL server",
		)
	}

	ctx, cancel := context.WithTimeout(r.ctx, 3*time.Second)
	defer cancel()

	if err := r.conn.Ping(ctx); err != nil {
		return repository.WrapErrorf(
			repository.ErrorFailedToConnect,
			err,
			"failed to ping to PostgreSQL server",
		)
	}

	return nil
}
