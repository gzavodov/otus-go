package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
	"github.com/jackc/pgx/v4"
)

//StatisticsRepository Storage interface for Statistics
type StatisticsRepository struct {
	BaseRepository
}

//NewStatisticsRepository creates new SQL Statistics Repository
func NewStatisticsRepository(ctx context.Context, dataSourceName string) repository.StatisticsRepository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &StatisticsRepository{BaseRepository{ctx: ctx, dataSourceName: dataSourceName}}
}

//Create creates new Statistics in databse
//If succseed ID field will be updated
func (r *StatisticsRepository) Create(m *model.Statistics) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	_, err := r.Execute(
		`INSERT INTO banner_statistics(banner_id, group_id, number_of_shows, number_of_clicks) VALUES($1, $2, $3, $4)
			ON CONFLICT(banner_id, group_id) DO UPDATE SET number_of_shows = $3, number_of_clicks = $4`,
		m.BannerID,
		m.GroupID,
		m.NumberOfShows,
		m.NumberOfClicks,
	)

	if err != nil {
		return repository.NewCreationError(err, "failed to execute insert query")
	}

	return nil
}

//Read reads Statistics from databse by ID
func (r *StatisticsRepository) Read(bannerID int64, groupID int64) (*model.Statistics, error) {
	if bannerID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if groupID <= 0 {
		return nil, repository.NewInvalidArgumentError("second parameter must be greater than zero")
	}

	row, err := r.QueryRow(
		`SELECT banner_id, group_id, number_of_shows, number_of_clicks FROM banner_statistics WHERE banner_id = $1 AND group_id = $2`,
		bannerID,
		groupID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}

	m := &model.Statistics{}
	if err := row.Scan(&m.BannerID, &m.GroupID, &m.NumberOfShows, &m.NumberOfClicks); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewNotFoundError("failed to find record with bannerID: %d and groupID: %d", bannerID, groupID)
		}

		return nil, repository.NewReadingError(err, "failed to fetch query result")
	}

	return m, nil
}

//Update modifies Statistics in databse
func (r *StatisticsRepository) Update(m *model.Statistics) error {
	if m == nil {
		return repository.NewInvalidArgumentError("first parameter must be not null pointer")
	}

	bannerID := m.BannerID
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("model bannerID must be greater than zero")
	}

	groupID := m.GroupID
	if groupID <= 0 {
		return repository.NewInvalidArgumentError("model groupID must be greater than zero")
	}

	result, err := r.Execute(
		`UPDATE banner_statistics SET number_of_shows = $1, number_of_clicks = $2 WHERE banner_id = $3 AND group_id = $4`,
		m.NumberOfShows,
		m.NumberOfClicks,
		bannerID,
		groupID,
	)

	if err != nil {
		return repository.NewUpdatingError(err, "failed to execute update query for record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with bannerID: %d and groupID: %d", bannerID, groupID)
	}
	return nil
}

//Delete removes Statistics from databse
func (r *StatisticsRepository) Delete(bannerID int64, groupID int64) error {
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if groupID <= 0 {
		return repository.NewInvalidArgumentError("second parameter must be greater than zero")
	}

	result, err := r.Execute(
		`DELETE FROM banner_statistics WHERE banner_id = $1 AND group_id = $2`,
		bannerID,
		groupID,
	)

	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with bannerID: %d and groupID: %d", bannerID, groupID)
	}
	return nil
}

//DeleteByBannerID removes all statistics associated with banner specified by bannerID
func (r *StatisticsRepository) DeleteByBannerID(bannerID int64) error {
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	_, err := r.Execute(`DELETE FROM banner_statistics WHERE banner_id = $1`, bannerID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with bannerID: %d", bannerID)
	}

	return nil
}

//DeleteByGroupID removes all statistivs associated with group specified by groupID
func (r *StatisticsRepository) DeleteByGroupID(groupID int64) error {
	if groupID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	_, err := r.Execute(`DELETE FROM banner_statistics WHERE group_id = $1`, groupID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute delete query for record with groupID: %d", groupID)
	}

	return nil
}

//GetBannerStatistics returns statistics associated with banner specified by bannerID
func (r *StatisticsRepository) GetBannerStatistics(bannerID int64) ([]*model.Statistics, error) {
	if bannerID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	rows, err := r.Query(
		`SELECT banner_id, group_id, number_of_shows, number_of_clicks FROM banner_statistics WHERE banner_id = $1`,
		bannerID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}
	defer rows.Close()

	list := make([]*model.Statistics, 0)
	for rows.Next() {
		m := &model.Statistics{}
		if err := rows.Scan(&m.BannerID, &m.GroupID, &m.NumberOfShows, &m.NumberOfClicks); err != nil {
			return nil, repository.NewReadingError(err, "failed to execute select query")
		}
		list = append(list, m)
	}
	return list, nil
}

//GetGroupStatistics returns statistics associated with slot specified by groupID
func (r *StatisticsRepository) GetGroupStatistics(groupID int64) ([]*model.Statistics, error) {
	if groupID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	rows, err := r.Query(
		`SELECT banner_id, group_id, number_of_shows, number_of_clicks FROM banner_statistics WHERE group_id = $1`,
		groupID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}
	defer rows.Close()

	list := make([]*model.Statistics, 0)
	for rows.Next() {
		m := &model.Statistics{}
		if err := rows.Scan(&m.BannerID, &m.GroupID, &m.NumberOfShows, &m.NumberOfClicks); err != nil {
			return nil, repository.NewReadingError(err, "failed to execute select query")
		}
		list = append(list, m)
	}
	return list, nil
}

//GetRotationStatistics returns rotation statistics for specified slot associated with specified group
func (r *StatisticsRepository) GetRotationStatistics(slotID int64, groupID int64) ([]*model.Statistics, error) {
	if slotID <= 0 {
		return nil, repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if groupID <= 0 {
		return nil, repository.NewInvalidArgumentError("second parameter must be greater than zero")
	}

	rows, err := r.Query(
		`SELECT b.banner_id, s.group_id, s.number_of_shows, s.number_of_clicks FROM banner_binding b 
			LEFT JOIN banner_statistics s ON b.banner_id = s.banner_id AND s.group_id = $1 WHERE b.slot_id = $2`,
		groupID,
		slotID,
	)

	if err != nil {
		return nil, repository.NewReadingError(err, "failed to execute select query")
	}
	defer rows.Close()

	list := make([]*model.Statistics, 0)
	for rows.Next() {
		m := &model.Statistics{}
		if err := rows.Scan(&m.BannerID, &m.GroupID, &m.NumberOfShows, &m.NumberOfClicks); err != nil {
			return nil, repository.NewReadingError(err, "failed to execute select query")
		}
		list = append(list, m)
	}
	return list, nil
}

//IncrementNumberOfShows increases Number Of Shows of statistics specified by bannerID and groupID
func (r *StatisticsRepository) IncrementNumberOfShows(bannerID int64, groupID int64) error {
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if groupID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	result, err := r.Execute(`UPDATE banner_statistics SET number_of_shows = (number_of_shows + 1) WHERE banner_id = $1 AND group_id = $2`, bannerID, groupID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute update query for record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	return nil
}

//IncrementNumberOfClicks increases Number Of Clicks of statistics specified by bannerID and groupID
func (r *StatisticsRepository) IncrementNumberOfClicks(bannerID int64, groupID int64) error {
	if bannerID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	if groupID <= 0 {
		return repository.NewInvalidArgumentError("first parameter must be greater than zero")
	}

	result, err := r.Execute(`UPDATE banner_statistics SET number_of_clicks = (number_of_clicks + 1) WHERE banner_id = $1 AND group_id = $2`, bannerID, groupID)
	if err != nil {
		return repository.NewDeletionError(err, "failed to execute update query for record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	if !result {
		return repository.NewNotFoundError("failed to find record with bannerID: %d and groupID: %d", bannerID, groupID)
	}

	return nil
}
