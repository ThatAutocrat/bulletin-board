package models

import "database/sql"

type RSVP struct {
	Going    int
	Maybe    int
	NotGoing int
	UserStatus string
}

func UpsertRSVP(db *sql.DB, postID, userID int, status string) error {
	_, err := db.Exec(`
		INSERT INTO rsvps (post_id, user_id, status) VALUES (?, ?, ?)
		ON CONFLICT(post_id, user_id) DO UPDATE SET status = excluded.status`,
		postID, userID, status)
	return err
}

func GetRSVPs(db *sql.DB, postID, userID int) (*RSVP, error) {
	r := &RSVP{}
	db.QueryRow(`SELECT COUNT(*) FROM rsvps WHERE post_id=? AND status='going'`, postID).Scan(&r.Going)
	db.QueryRow(`SELECT COUNT(*) FROM rsvps WHERE post_id=? AND status='maybe'`, postID).Scan(&r.Maybe)
	db.QueryRow(`SELECT COUNT(*) FROM rsvps WHERE post_id=? AND status='not_going'`, postID).Scan(&r.NotGoing)
	if userID > 0 {
		var s sql.NullString
		db.QueryRow(`SELECT status FROM rsvps WHERE post_id=? AND user_id=?`, postID, userID).Scan(&s)
		r.UserStatus = s.String
	}
	return r, nil
}
