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

//QueryRow is a wrapper over pgx QueryRow
func (r *BaseRepository) QueryRow(query string, params ...interface{}) (pgx.Row, error) {
	if err := r.Connect(); err != nil {
		return nil, err
	}

	return r.conn.QueryRow(r.ctx, query, params...), nil
}

//Query is a wrapper over pgx Query
func (r *BaseRepository) Query(query string, params ...interface{}) (pgx.Rows, error) {
	if err := r.Connect(); err != nil {
		return nil, err
	}

	return r.conn.Query(r.ctx, query, params...)
}

//Execute is a wrapper over pgx Exec
func (r *BaseRepository) Execute(query string, params ...interface{}) (bool, error) {
	if err := r.Connect(); err != nil {
		return false, err
	}

	res, err := r.conn.Exec(r.ctx, query, params...)
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, nil
}
