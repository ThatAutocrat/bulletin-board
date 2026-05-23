package handlers

import (
	"net/http"

	"github.com/ThatAutocrat/bulletin-board/internal/middleware"
	"github.com/ThatAutocrat/bulletin-board/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) RegisterGet(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "register.html", nil)
}

func (a *App) RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.render(w, r, "register.html", &TemplateData{Error: "Something went wrong"})
		return
	}
	_, err = models.CreateUser(a.DB, username, email, string(hash))
	if err != nil {
		a.render(w, r, "register.html", &TemplateData{Error: "Username or email already taken"})
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (a *App) LoginGet(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "login.html", nil)
}

func (a *App) LoginPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(a.DB, email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		a.render(w, r, "login.html", &TemplateData{Error: "Invalid email or password"})
		return
	}
	sessionID, err := models.CreateSession(a.DB, user.ID)
	if err != nil {
		a.render(w, r, "login.html", &TemplateData{Error: "Could not create session"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   7 * 24 * 3600,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		models.DeleteSession(a.DB, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session", MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	user, err := models.GetUserByUsername(a.DB, username)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	posts, _ := models.GetPostsByUser(a.DB, user.ID)
	_ = middleware.GetUserID(r)
	a.render(w, r, "profile.html", &TemplateData{Posts: posts, Post: &models.Post{Username: user.Username}})
}
