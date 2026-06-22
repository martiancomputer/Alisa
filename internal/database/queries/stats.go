package queries

import (
	"database/sql"
	"time"
)

type StatsRepository struct {
	DB *sql.DB
}

// IncrementHourlyMetric increments or instantiates an aggregate counter for a metrics key.
func (r *StatsRepository) IncrementHourlyMetric(metricType, targetID string) error {
	// Truncate timestamp to current hour block boundaries
	currentHour := time.Now().Truncate(time.Hour).Unix()

	query := `
		INSERT INTO statistics_hourly (timestamp_hour, metric_type, target_id, value_count)
		VALUES (?, ?, ?, 1)
		ON CONFLICT(timestamp_hour, metric_type, target_id) DO UPDATE SET
		value_count = value_count + 1;
	`
	_, err := r.DB.Exec(query, currentHour, metricType, targetID)
	return err
}