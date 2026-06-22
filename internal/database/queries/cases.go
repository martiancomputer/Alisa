package queries

import (
	"database/sql"
	"time"

	"github.com/martiancomputer/Alisa/internal/models"
)

type CaseRepository struct {
	DB *sql.DB
}

// LogCase writes a new infraction record to the local instance.
func (r *CaseRepository) LogCase(c *models.Case) (int64, error) {
	query := `
		INSERT INTO cases (user_id, moderator_id, action_type, reason, duration_seconds, timestamp)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	res, err := r.DB.Exec(query, c.UserID, c.ModeratorID, c.ActionType, c.Reason, c.DurationSeconds, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetUserHistory reads all archived infractions tied to a unique snowflake string.
func (r *CaseRepository) GetUserHistory(userID string) ([]models.Case, error) {
	query := `
		SELECT case_id, user_id, moderator_id, action_type, reason, duration_seconds, timestamp
		FROM cases
		WHERE user_id = ?
		ORDER BY timestamp DESC;
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.Case
	for rows.Next() {
		var c models.Case
		err := rows.Scan(&c.ID, &c.UserID, &c.ModeratorID, &c.ActionType, &c.Reason, &c.DurationSeconds, &c.Timestamp)
		if err != nil {
			return nil, err
		}
		history = append(history, c)
	}
	return history, nil
}