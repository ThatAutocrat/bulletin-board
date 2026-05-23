package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ThatAutocrat/bulletin-board/internal/middleware"
	"github.com/ThatAutocrat/bulletin-board/internal/models"
)

type App struct {
	DB *sql.DB
}

type TemplateData struct {
	User          *models.User
	Posts         []models.Post
	Post          *models.Post
	Comments      []models.Comment
	RSVP          *models.RSVP
	UpcomingEvents []models.Post
	Category      string
	Error         string
	Success       string
}

func (a *App) render(w http.ResponseWriter, r *http.Request, page string, data *TemplateData) {
	if data == nil {
		data = &TemplateData{}
	}
	userID := middleware.GetUserID(r)
	if userID > 0 {
		u, err := models.GetUserByID(a.DB, userID)
		if err == nil {
			data.User = u
		}
	}
	events, _ := models.GetUpcomingEvents(a.DB)
	data.UpcomingEvents = events

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string { return t.Format("Jan 2, 2006") },
		"formatDateTime": func(t time.Time) string { return t.Format("Jan 2, 2006 3:04 PM") },
		"categoryLabel": func(c string) string {
			labels := map[string]string{
				"announcement": "📢 Announcement",
				"lost_found":   "🔍 Lost & Found",
				"event":        "📅 Event",
				"general":      "💬 General",
			}
			if l, ok := labels[c]; ok { return l }
			return c
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "pages", page),
	)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), 500)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Render error: "+err.Error(), 500)
	}
}
