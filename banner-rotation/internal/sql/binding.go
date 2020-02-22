package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
	"github.com/jackc/pgx/v4"
)

//BindingRepository Storage interface for Banner Bindings
type BindingRepository struct {
	BaseRepository
}

//NewBindingRepository creates new SQL Binding Repository
func NewBindingRepository(ctx context.Context, dataSourceName string) repository.BindingRepository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &BindingRepository{BaseRepository{ctx: ctx, dataSourceName: dataSourceName}}
}

//Create creates new Binding in databse
//If succseed ID field will be updated
func (r *BindingRepository) Create(m *model.Binding) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	//For operationg in concurrecy mode we are using upsert operation. The fake update on conflict is required for return actual record ID.
	row, err := r.QueryRow(
		`INSERT INTO banner_binding AS b(banner_id, slot_id) VALUES($1, $2) 
			ON CONFLICT(banner_id, slot_id) 
			DO UPDATE SET banner_id = b.banner_id, slot_id = b.slot_id 
		RETURNING b.id`,
		m.BannerID,
		m.SlotID,
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

//Read reads Banner Binding from databse by ID
func (r *BindingRepository) Read(ID int64) (*model.Binding, error) {
	if ID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(
		`SELECT id, banner_id, slot_id FROM banner_binding WHERE id = $1`,
		ID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}

	m := &model.Binding{}
	if err := row.Scan(&m.ID, &m.BannerID, &m.SlotID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewNotFoundError("failed to find record with ID: %d", ID)
		}

		return nil, repository.NewReadingError(err, "failed to fetch query result")
	}

	return m, nil
}

//Update modifies Binding in databse
func (r *BindingRepository) Update(m *model.Binding) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	ID := m.ID
	if ID <= 0 {
		return repository.NewInvalidArgumentError("model ID must be greater than zero")
	}

	bannerID := m.BannerID
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("model bannerID must be greater than zero")
	}

	slotID := m.SlotID
	if slotID <= 0 {
		return repository.NewInvalidArgumentError("model slotID must be greater than zero")
	}

	result, err := r.Execute(
		`UPDATE banner_binding SET banner_id = $1, slot_id = $2 WHERE id = $3`,
		bannerID,
		slotID,
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

//Delete removes Binding from databse
func (r *BindingRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	result, err := r.Execute(`DELETE FROM banner_binding WHERE id = $1`, ID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with ID: %d", ID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with ID: %d", ID)
	}
	return nil
}

//DeleteByBannerID removes all bindings associated with banner specified by bannerID
func (r *BindingRepository) DeleteByBannerID(bannerID int64) error {
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	_, err := r.Execute(`DELETE FROM banner_binding WHERE banner_id = $1`, bannerID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with bannerID: %d", bannerID)
	}

	return nil
}

//DeleteBySlotID removes all bindings associated with slot specified by slotID
func (r *BindingRepository) DeleteBySlotID(slotID int64) error {
	if slotID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	_, err := r.Execute(`DELETE FROM banner_binding WHERE slot_id = $1`, slotID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with slotID: %d", slotID)
	}

	return nil
}

//DeleteByModel removes all bindings associated with banner and slot specified by item
func (r *BindingRepository) DeleteByModel(m *model.Binding) error {
	bannerID := m.BannerID
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("model bannerID must be greater than zero")
	}

	slotID := m.SlotID
	if slotID <= 0 {
		return repository.NewInvalidArgumentError("model slotID must be greater than zero")
	}

	_, err := r.Execute(`DELETE FROM banner_binding WHERE banner_id = $1 AND slot_id = $2`, bannerID, slotID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with slotID: %d", slotID)
	}

	return nil
}

//GetBinding returns binding associated with specified banner and slot
func (r *BindingRepository) GetBinding(bannerID int64, slotID int64) (*model.Binding, error) {
	if bannerID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if slotID <= 0 {
		return nil, repository.NewInvalidArgumentError("second parameter must be greater than zero")
	}

	row, err := r.QueryRow(
		`SELECT id, banner_id, slot_id FROM banner_binding WHERE banner_id = $1 AND slot_id = $2`,
		bannerID,
		slotID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}

	m := &model.Binding{}
	if err := row.Scan(&m.ID, &m.BannerID, &m.SlotID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, repository.NewReadingError(err, "failed to fetch query result")
	}

	return m, nil
}

//GetBannerBindings returns bindings associated with banner specified by bannerID
func (r *BindingRepository) GetBannerBindings(bannerID int64) ([]*model.Binding, error) {
	if bannerID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	rows, err := r.Query(
		`SELECT id, banner_id, slot_id FROM banner_binding WHERE banner_id = $1`,
		bannerID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}
	defer rows.Close()

	list := make([]*model.Binding, 0)
	for rows.Next() {
		m := &model.Binding{}
		if err := rows.Scan(&m.ID, &m.BannerID, &m.SlotID); err != nil {
			return nil, repository.NewReadingError(err, "failed to execute select query")
		}
		list = append(list, m)
	}
	return list, nil
}

//GetSlotBindings returns bindings associated with slot specified by slotID
func (r *BindingRepository) GetSlotBindings(slotID int64) ([]*model.Binding, error) {
	if slotID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	rows, err := r.Query(
		`SELECT id, banner_id, slot_id FROM banner_binding WHERE slot_id = $1`,
		slotID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}
	defer rows.Close()

	list := make([]*model.Binding, 0)
	for rows.Next() {
		m := &model.Binding{}
		if err := rows.Scan(&m.ID, &m.BannerID, &m.SlotID); err != nil {
			return nil, repository.NewReadingError(err, "failed to execute select query")
		}
		list = append(list, m)
	}
	return list, nil
}

//IsExists check if repository contains banner with specified ID
func (r *BindingRepository) IsExists(ID int64) (bool, error) {
	if ID <= 0 {
		return false, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	row, err := r.QueryRow(`SELECT 'x' FROM banner_binding WHERE id = $1`, ID)
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
