package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID          int
	UserID      int
	Username    string
	Category    string
	Title       string
	Body        string
	LocationTag string
	Upvotes     int
	Status      string
	EventDate   *time.Time
	ItemType    string
	CreatedAt   time.Time
	CommentCount int
}

func CreatePost(db *sql.DB, p *Post) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO posts (user_id, category, title, body, location_tag, event_date, item_type)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.UserID, p.Category, p.Title, p.Body, p.LocationTag, p.EventDate, p.ItemType,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetPosts(db *sql.DB, category string) ([]Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.category, p.title, p.body,
		       p.location_tag, p.upvotes, p.status, p.event_date, p.item_type, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comment_count
		FROM posts p JOIN users u ON u.id = p.user_id`
	args := []any{}
	if category != "" {
		query += ` WHERE p.category = ?`
		args = append(args, category)
	}
	query += ` ORDER BY p.created_at DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var loc, itemType sql.NullString
		var evDate sql.NullTime
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Category, &p.Title, &p.Body,
			&loc, &p.Upvotes, &p.Status, &evDate, &itemType, &p.CreatedAt, &p.CommentCount); err != nil {
			return nil, err
		}
		p.LocationTag = loc.String
		p.ItemType = itemType.String
		if evDate.Valid {
			p.EventDate = &evDate.Time
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, id int) (*Post, error) {
	p := &Post{}
	var loc, itemType sql.NullString
	var evDate sql.NullTime
	err := db.QueryRow(`
		SELECT p.id, p.user_id, u.username, p.category, p.title, p.body,
		       p.location_tag, p.upvotes, p.status, p.event_date, p.item_type, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
		FROM posts p JOIN users u ON u.id = p.user_id WHERE p.id = ?`, id,
	).Scan(&p.ID, &p.UserID, &p.Username, &p.Category, &p.Title, &p.Body,
		&loc, &p.Upvotes, &p.Status, &evDate, &itemType, &p.CreatedAt, &p.CommentCount)
	if err != nil {
		return nil, err
	}
	p.LocationTag = loc.String
	p.ItemType = itemType.String
	if evDate.Valid {
		p.EventDate = &evDate.Time
	}
	return p, nil
}

func UpvotePost(db *sql.DB, id int) error {
	_, err := db.Exec(`UPDATE posts SET upvotes = upvotes + 1 WHERE id = ?`, id)
	return err
}

func ToggleStatus(db *sql.DB, id int) error {
	_, err := db.Exec(`
		UPDATE posts SET status = CASE WHEN status='open' THEN 'resolved' ELSE 'open' END
		WHERE id = ?`, id)
	return err
}

func GetPostsByUser(db *sql.DB, userID int) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, u.username, p.category, p.title, p.body,
		       p.location_tag, p.upvotes, p.status, p.event_date, p.item_type, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
		FROM posts p JOIN users u ON u.id = p.user_id
		WHERE p.user_id = ? ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		var loc, itemType sql.NullString
		var evDate sql.NullTime
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Category, &p.Title, &p.Body,
			&loc, &p.Upvotes, &p.Status, &evDate, &itemType, &p.CreatedAt, &p.CommentCount); err != nil {
			return nil, err
		}
		p.LocationTag = loc.String
		p.ItemType = itemType.String
		if evDate.Valid {
			p.EventDate = &evDate.Time
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetUpcomingEvents(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, u.username, p.category, p.title, p.body,
		       p.location_tag, p.upvotes, p.status, p.event_date, p.item_type, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
		FROM posts p JOIN users u ON u.id = p.user_id
		WHERE p.category = 'event' AND p.event_date >= datetime('now')
		ORDER BY p.event_date ASC LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		var loc, itemType sql.NullString
		var evDate sql.NullTime
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Category, &p.Title, &p.Body,
			&loc, &p.Upvotes, &p.Status, &evDate, &itemType, &p.CreatedAt, &p.CommentCount); err != nil {
			return nil, err
		}
		p.LocationTag = loc.String
		p.ItemType = itemType.String
		if evDate.Valid {
			p.EventDate = &evDate.Time
		}
		posts = append(posts, p)
	}
	return posts, nil
}
