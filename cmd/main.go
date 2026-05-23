package main

import (
	"log"
	"net/http"
	"os"

	"bulletin/internal/handlers"
	"bulletin/internal/middleware"
	"bulletin/internal/models"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./bulletin.db"
	}
	db, err := models.InitDB(dbPath)
	if err != nil {
		log.Fatal("DB init error:", err)
	}
	defer db.Close()

	app := &handlers.App{DB: db}
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.Auth(db))

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", app.Home)

	// Auth
	r.Get("/auth/register", app.RegisterGet)
	r.Post("/auth/register", app.RegisterPost)
	r.Get("/auth/login", app.LoginGet)
	r.Post("/auth/login", app.LoginPost)
	r.Post("/auth/logout", app.Logout)
	r.Get("/profile/{username}", app.Profile)

	// Posts
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/posts/new", app.NewPostGet)
		r.Post("/posts", app.NewPostPost)
		r.Post("/posts/{id}/comment", app.AddComment)
		r.Post("/posts/{id}/upvote", app.Upvote)
		r.Post("/posts/{id}/toggle-status", app.ToggleStatus)
		r.Post("/posts/{id}/rsvp", app.RSVP)
	})
	r.Get("/posts/{id}", app.ViewPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
