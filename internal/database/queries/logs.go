package queries

import (
	"database/sql"
	"time"
)

type LogRepository struct {
	DB *sql.DB
}

// WriteAuditRecord pushes a structured server modification event to the persistent stream.
func (r *LogRepository) WriteAuditRecord(eventType, userID, moderatorID, channelID, payload string) error {
	query := `
		INSERT INTO audit_log (event_type, user_id, moderator_id, channel_id, timestamp, payload)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err := r.DB.Exec(query, eventType, userID, moderatorID, channelID, time.Now().Unix(), payload)
	return err
}

// FetchAuditFeed retrieves paginated telemetry matching specific criteria.
func (r *LogRepository) FetchAuditFeed(limit, offset int) (*sql.Rows, error) {
	query := `
		SELECT log_id, event_type, user_id, moderator_id, channel_id, timestamp, payload
		FROM audit_log
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;
	`
	return r.DB.Query(query, limit, offset)
}