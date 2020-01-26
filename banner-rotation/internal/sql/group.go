package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
	"github.com/jackc/pgx/v4"
)

//GroupRepository Storage interface for Banner Group
type GroupRepository struct {
	BaseRepository
}

//NewGroupRepository creates new SQL Group Repository
func NewGroupRepository(ctx context.Context, dataSourceName string) *GroupRepository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &GroupRepository{BaseRepository{ctx: ctx, dataSourceName: dataSourceName}}
}

//Create creates new Banner Group in databse
//If succseed ID field will be updated
func (r *GroupRepository) Create(m *model.Group) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	if err := r.Connect(); err != nil {
		return err
	}

	row := r.conn.QueryRow(
		r.ctx,
		`INSERT INTO group(caption) VALUES($1) RETURNING id`,
		m.Caption,
	)

	//sql.Row.Scan will close underlying sql.Row before exit
	if err := row.Scan(&m.ID); err != nil {
		return repository.NewCreationError(err, "failed to execute insert query")
	}

	return nil
}

//Read reads Banner Group from databse by ID
func (r *GroupRepository) Read(ID int64) (*model.Group, error) {
	if ID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if err := r.Connect(); err != nil {
		return nil, err
	}

	row := r.conn.QueryRow(
		r.ctx,
		`SELECT id, caption FROM group WHERE id = $1`,
		ID,
	)

	m := &model.Group{}
	err := row.Scan(&m.ID, &m.Caption)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewNotFoundError("failed to find record with ID: %d", ID)
		}

		return nil, repository.NewReadingError(err, "failed to read record with ID: %d", ID)
	}

	return m, nil
}

//Update modifies Banner Group in databse
func (r *GroupRepository) Update(m *model.Group) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	ID := m.ID
	if ID <= 0 {
		return repository.NewInvalidArgumentError("model ID must be greater than zero")
	}

	if err := r.Connect(); err != nil {
		return err
	}

	res, err := r.conn.Exec(
		r.ctx,
		`UPDATE group SET caption = $1 WHERE id = $2`,
		m.Caption,
		ID,
	)

	if err != nil {
		return repository.NewUpdatingError(err, "failed to execute update query for record with ID: %d", ID)
	}

	if res.RowsAffected() == 0 {
		return repository.NewNotFoundError("failed to find record with ID: %d", ID)
	}
	return nil
}

//Delete removes Banner Group from databse
func (r *GroupRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if err := r.Connect(); err != nil {
		return err
	}

	res, err := r.conn.Exec(r.ctx, `DELETE FROM group WHERE id = $1`, ID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with ID: %d", ID)
	}

	if res.RowsAffected() == 0 {
		return repository.NewNotFoundError("failed to find record with ID: %d", ID)
	}
	return nil
}
