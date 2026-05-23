package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Username  string
	Body      string
	CreatedAt time.Time
}

func CreateComment(db *sql.DB, postID, userID int, body string) error {
	_, err := db.Exec(
		`INSERT INTO comments (post_id, user_id, body) VALUES (?, ?, ?)`,
		postID, userID, body,
	)
	return err
}

func GetCommentsByPost(db *sql.DB, postID int) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.body, c.created_at
		FROM comments c JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ? ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
