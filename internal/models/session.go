package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID int) (string, error) {
	id := uuid.New().String()
	expires := time.Now().Add(7 * 24 * time.Hour)
	_, err := db.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		id, userID, expires,
	)
	return id, err
}

func GetUserIDBySession(db *sql.DB, sessionID string) (int, error) {
	var userID int
	err := db.QueryRow(
		`SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`, sessionID,
	).Scan(&userID)
	return userID, err
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}
