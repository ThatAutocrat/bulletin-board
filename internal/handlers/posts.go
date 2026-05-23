package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ThatAutocrat/bulletin-board/internal/middleware"
	"github.com/ThatAutocrat/bulletin-board/internal/models"
)

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	posts, err := models.GetPosts(a.DB, category)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	a.render(w, r, "index.html", &TemplateData{Posts: posts, Category: category})
}

func (a *App) NewPostGet(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "new_post.html", nil)
}

func (a *App) NewPostPost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	p := &models.Post{
		UserID:      userID,
		Category:    r.FormValue("category"),
		Title:       r.FormValue("title"),
		Body:        r.FormValue("body"),
		LocationTag: r.FormValue("location_tag"),
		ItemType:    r.FormValue("item_type"),
	}
	if ed := r.FormValue("event_date"); ed != "" {
		t, err := time.Parse("2006-01-02T15:04", ed)
		if err == nil {
			p.EventDate = &t
		}
	}
	id, err := models.CreatePost(a.DB, p)
	if err != nil {
		a.render(w, r, "new_post.html", &TemplateData{Error: "Could not create post"})
		return
	}
	http.Redirect(w, r, "/posts/"+strconv.FormatInt(id, 10), http.StatusFound)
}

func (a *App) ViewPost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	post, err := models.GetPostByID(a.DB, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comments, _ := models.GetCommentsByPost(a.DB, id)
	userID := middleware.GetUserID(r)
	rsvp, _ := models.GetRSVPs(a.DB, id, userID)
	a.render(w, r, "post.html", &TemplateData{Post: post, Comments: comments, RSVP: rsvp})
}

func (a *App) AddComment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userID := middleware.GetUserID(r)
	body := r.FormValue("body")
	if body != "" {
		models.CreateComment(a.DB, id, userID, body)
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (a *App) Upvote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	models.UpvotePost(a.DB, id)
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (a *App) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userID := middleware.GetUserID(r)
	post, err := models.GetPostByID(a.DB, id)
	if err != nil || post.UserID != userID {
		http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
		return
	}
	models.ToggleStatus(a.DB, id)
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (a *App) RSVP(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userID := middleware.GetUserID(r)
	status := r.FormValue("status")
	if status != "" {
		models.UpsertRSVP(a.DB, id, userID, status)
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}
