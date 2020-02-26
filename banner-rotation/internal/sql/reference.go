package sql

import (
	"fmt"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
	"github.com/jackc/pgx/v4"
)

type ReferenceRepository struct {
	TableName string
	BaseRepository
}

func (r *ReferenceRepository) DoCreate(m model.Reference) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	row, err := r.QueryRow(
		fmt.Sprintf(`INSERT INTO %s(caption) VALUES($1) RETURNING id`, r.TableName),
		m.GetCaption(),
	)

	if err != nil {
		return repository.NewCreationError(err, "failed to execute insert query")
	}

	//sql.Row.Scan will close underlying sql.Row before exit
	var id int64

	if err := row.Scan(&id); err != nil {
		return repository.NewCreationError(err, "failed to fetch query result")
	}

	m.SetID(id)
	return nil
}

func (r *ReferenceRepository) DoRead(id int64, m model.Reference) error {
	if id <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(
		fmt.Sprintf(`SELECT id, caption FROM %s WHERE id = $1`, r.TableName),
		id,
	)

	if err != nil {
		return repository.NewReadingError(err, "failed to execute select query")
	}

	var caption string

	if err := row.Scan(&id, &caption); err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewNotFoundError("failed to find record with id: %d", id)
		}

		return repository.NewReadingError(err, "failed to fetch query result")
	}

	m.SetID(id)
	m.SetCaption(caption)

	return nil
}

func (r *ReferenceRepository) DoUpdate(m model.Reference) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	ID := m.GetID()
	if ID <= 0 {
		return repository.NewInvalidArgumentError("model ID must be greater than zero")
	}

	result, err := r.Execute(
		fmt.Sprintf(`UPDATE %s SET caption = $1 WHERE id = $2`, r.TableName),
		m.GetCaption(),
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

func (r *ReferenceRepository) DoDelete(id int64) error {
	if id <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	result, err := r.Execute(fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.TableName), id)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with id: %d", id)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with id: %d", id)
	}
	return nil
}

func (r *ReferenceRepository) CheckIfExists(id int64) (bool, error) {
	if id <= 0 {
		return false, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(fmt.Sprintf(`SELECT 'x' FROM %s WHERE id = $1`, r.TableName), id)
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

func (r *ReferenceRepository) DoGetByCaption(caption string, m model.Reference) error {
	if caption == "" {
		return repository.NewInvalidArgumentError("first parameter must be not empty string")
	}

	row, err := r.QueryRow(
		fmt.Sprintf(`SELECT id, caption FROM %s WHERE caption = $1 ORDER BY id DESC LIMIT 1`, r.TableName),
		caption,
	)

	if err != nil {
		return repository.NewReadingError(err, "failed to execute select query")
	}

	var id int64

	if err := row.Scan(&id, &caption); err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewNotFoundError("failed to find record with caption: %s", caption)
		}

		return repository.NewReadingError(err, "failed to fetch query result")
	}

	m.SetID(id)
	m.SetCaption(caption)

	return nil
}
