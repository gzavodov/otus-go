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
func NewGroupRepository(ctx context.Context, dataSourceName string) repository.GroupRepository {
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

	row, err := r.QueryRow(
		`INSERT INTO banner_group(caption) VALUES($1) RETURNING id`,
		m.Caption,
	)

	if err != nil {
		return repository.NewCreationError(err, "failed to execute insert query")
	}

	//sql.Row.Scan will close underlying sql.Row before exit
	if err := row.Scan(&m.ID); err != nil {
		return repository.NewCreationError(err, "failed to fetch query result")
	}

	return nil
}

//Read reads Banner Group from databse by ID
func (r *GroupRepository) Read(ID int64) (*model.Group, error) {
	if ID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(
		`SELECT id, caption FROM banner_group WHERE id = $1`,
		ID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}

	m := &model.Group{}
	if err := row.Scan(&m.ID, &m.Caption); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewNotFoundError("failed to find record with ID: %d", ID)
		}

		return nil, repository.NewReadingError(err, "failed to fetch query result")
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

	result, err := r.Execute(
		`UPDATE banner_group SET caption = $1 WHERE id = $2`,
		m.Caption,
		ID,
	)

	if err != nil {
		return repository.NewUpdatingError(err, "failed to execute update query for record with ID: %d", ID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with ID: %d", ID)
	}
	return nil
}

//Delete removes Banner Group from databse
func (r *GroupRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	result, err := r.Execute(`DELETE FROM banner_group WHERE id = $1`, ID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with ID: %d", ID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with ID: %d", ID)
	}
	return nil
}

//IsExists check if repository contains Banner Group with specified ID
func (r *GroupRepository) IsExists(ID int64) (bool, error) {
	if ID <= 0 {
		return false, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(`SELECT 'x' FROM banner_group WHERE id = $1`, ID)
	if err != nil {
		return false, repository.NewReadingError(err, "failed to execute select query")
	}

	s := ""
	if err := row.Scan(&s); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, repository.NewReadingError(err, "failed to fetch query result")
	}

	return true, nil
}

//GetByCaption returns Banner Group bt specified caption
func (r *GroupRepository) GetByCaption(caption string) (*model.Group, error) {
	if caption == "" {
		return nil, repository.NewInvalidArgumentError("first parameter must be not empty string")
	}

	row, err := r.QueryRow(
		`SELECT id, caption FROM banner_group WHERE caption = $1 ORDER BY id DESC LIMIT 1`,
		caption,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}

	m := &model.Group{}
	if err := row.Scan(&m.ID, &m.Caption); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewNotFoundError("failed to find record with caption: %s", caption)
		}

		return nil, repository.NewReadingError(err, "failed to fetch query result")
	}

	return m, nil
}
