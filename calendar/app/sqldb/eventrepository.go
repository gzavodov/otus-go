package sqldb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
)

//NewEventRepository creates new SQL EventRepository
func NewEventRepository(ctx context.Context, dataSourceName string) *EventRepository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &EventRepository{ctx: ctx, dataSourceName: dataSourceName}
}

//EventRepository SQL implementation of EventRepository interface
type EventRepository struct {
	dataSourceName string
	conn           *pgx.Conn
	ctx            context.Context
}

//Connect try to connect to PostgreSQL server
func (r *EventRepository) Connect() error {
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

//Create add Calendar Event in repository
//If succseed ID field will be updated
func (r *EventRepository) Create(m *model.Event) error {
	if m == nil {
		return repository.NewError(
			repository.ErrorInvalidArgument,
			"first parameter must be not null pointer to event",
		)
	}

	if err := r.Connect(); err != nil {
		return err
	}

	now := time.Now().UTC()
	row := r.conn.QueryRow(
		r.ctx,
		`INSERT INTO event(title, 
			description, location, 
			start_time, end_time, 
			user_id, calendar_id, 
			created, last_updated
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		m.Title,
		m.Description,
		m.Location,
		m.StartTime.UTC(),
		m.EndTime.UTC(),
		m.UserID,
		m.CalendarID,
		now,
		now,
	)

	//sql.Row.Scan will close underlying sql.Row before exit
	if err := row.Scan(&m.ID); err != nil {
		return repository.WrapErrorf(
			repository.ErrorDataCreationFailure,
			err,
			"failed to execute insert query",
		)
	}

	return nil
}

//IsExists check if repository contains Calendar event with specified ID
func (r *EventRepository) IsExists(ID int64) (bool, error) {
	if ID <= 0 {
		return false,
			repository.NewError(
				repository.ErrorInvalidArgument,
				"first parameter must be greater than zero",
			)
	}

	if err := r.Connect(); err != nil {
		return false, err
	}

	row := r.conn.QueryRow(r.ctx, `SELECT 'x' FROM event WHERE id = $1`, ID)

	s := ""
	if err := row.Scan(&s); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false,
			repository.WrapErrorf(
				repository.ErrorDataRetrievalFailure,
				err,
				"failed to execute select query for record with ID: %d",
				ID,
			)
	}

	return true, nil
}

//Read get Calendar Event from repository by ID
func (r *EventRepository) Read(ID int64) (*model.Event, error) {
	if ID <= 0 {
		return nil,
			repository.NewError(
				repository.ErrorInvalidArgument,
				"first parameter must be greater than zero",
			)
	}

	if err := r.Connect(); err != nil {
		return nil, err
	}

	row := r.conn.QueryRow(
		r.ctx,
		`SELECT id, title, description, location, start_time, end_time, user_id, calendar_id FROM event WHERE id = $1`,
		ID,
	)

	m := &model.Event{}
	err := row.Scan(&m.ID, &m.Title, &m.Description, &m.Location, &m.StartTime, &m.EndTime, &m.UserID, &m.CalendarID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil,
				repository.NewError(
					repository.ErrorNotFound,
					fmt.Sprintf("failed to find record with ID: %d", ID),
				)
		}
		return nil,
			repository.WrapErrorf(
				repository.ErrorDataRetrievalFailure,
				err,
				"failed to read record with ID: %d",
				ID,
			)
	}

	m.StartTime = m.StartTime.UTC()
	m.EndTime = m.EndTime.UTC()

	return m, nil
}

//ReadList get Calendar Events by interval specified by from and to params
func (r *EventRepository) ReadList(userID int64, from time.Time, to time.Time) ([]*model.Event, error) {
	if err := r.Connect(); err != nil {
		return nil, err
	}

	from = from.UTC()
	to = to.UTC()

	conditions := make([]string, 0, 3)
	params := make([]interface{}, 0, 3)
	if userID > 0 {
		params = append(params, userID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(params)))
	}

	if !from.IsZero() {
		params = append(params, from)
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", len(params)))
	}

	if !to.IsZero() {
		params = append(params, to)
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", len(params)))
	}

	query := `SELECT id, title, description, location, start_time, end_time, user_id, calendar_id FROM event`
	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s ORDER BY id ASC", query, strings.Join(conditions, " AND "))
	} else {
		query = fmt.Sprintf("%s ORDER BY id ASC", query)
	}

	rows, err := r.conn.Query(r.ctx, query, params...)
	if err != nil {
		return nil,
			repository.WrapError(
				repository.ErrorDataRetrievalFailure,
				err,
			)
	}
	defer rows.Close()

	list := make([]*model.Event, 0, 10)
	for rows.Next() {
		m := &model.Event{}
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Location, &m.StartTime, &m.EndTime, &m.UserID, &m.CalendarID)
		if err != nil {
			return nil,
				repository.WrapError(
					repository.ErrorDataRetrievalFailure,
					err,
				)
		}
		m.StartTime = m.StartTime.UTC()
		m.EndTime = m.EndTime.UTC()

		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil,
			repository.WrapError(
				repository.ErrorDataRetrievalFailure,
				err,
			)
	}

	return list, nil
}

//ReadAll get all Calendar Events from repository
func (r *EventRepository) ReadAll() ([]*model.Event, error) {
	return r.ReadList(0, time.Time{}, time.Time{})
}

//Update modifies Calendar Event in repository
func (r *EventRepository) Update(m *model.Event) error {
	if m == nil {
		return repository.NewError(
			repository.ErrorInvalidArgument,
			"first parameter must be not null pointer to event",
		)
	}

	ID := m.ID
	if ID <= 0 {
		return repository.NewError(
			repository.ErrorInvalidArgument,
			"model ID must be greater than zero",
		)
	}

	if err := r.Connect(); err != nil {
		return err
	}

	res, err := r.conn.Exec(
		r.ctx,
		`UPDATE event SET
			title = $1, description = $2, location = $3, 
			start_time = $4, end_time = $5, 
			user_id = $6, calendar_id = $7, 
			last_updated = $8
		WHERE id = $9`,
		m.Title,
		m.Description,
		m.Location,
		m.StartTime.UTC(),
		m.EndTime.UTC(),
		m.UserID,
		m.CalendarID,
		time.Now().UTC(),
		ID,
	)

	if err != nil {
		return repository.WrapErrorf(
			repository.ErrorDataModificationFailure,
			err,
			"failed to execute update query for record with ID: %d",
			ID,
		)
	}

	if res.RowsAffected() == 0 {
		return repository.NewErrorf(
			repository.ErrorNotFound,
			"failed to find record with ID: %d",
			ID,
		)
	}
	return nil
}

//Delete removes Calendar Event from repository by ID
func (r *EventRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewError(
			repository.ErrorInvalidArgument,
			"first parameter must be greater than zero",
		)
	}

	if err := r.Connect(); err != nil {
		return err
	}

	res, err := r.conn.Exec(r.ctx, `DELETE FROM event WHERE id = $1`, ID)
	if err != nil {
		return repository.WrapErrorf(
			repository.ErrorDataDeletionFailure,
			err,
			"failed to execute delete query for record with ID: %d",
			ID,
		)
	}

	if res.RowsAffected() == 0 {
		return repository.NewErrorf(
			repository.ErrorNotFound,
			"failed to find record with ID: %d",
			ID,
		)
	}
	return nil
}

//GetTotalCount returns overall amouunt of calendar events in repository
func (r *EventRepository) GetTotalCount() (int64, error) {
	if err := r.Connect(); err != nil {
		return 0, err
	}

	var result int64
	row := r.conn.QueryRow(r.ctx, `SELECT COUNT(*) FROM event`)
	if err := row.Scan(&result); err != nil {
		return 0,
			repository.WrapErrorf(
				repository.ErrorDataRetrievalFailure,
				err,
				"failed to execute select query",
			)
	}

	return result, nil
}

//Purge removes all Calendar records from repository
func (r *EventRepository) Purge() error {
	if err := r.Connect(); err != nil {
		return err
	}

	if _, err := r.conn.Exec(r.ctx, `DELETE FROM event`); err != nil {
		return repository.WrapErrorf(
			repository.ErrorDataRetrievalFailure,
			err,
			"failed to execute delete query",
		)
	}

	return nil
}
