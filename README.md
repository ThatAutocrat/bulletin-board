# 📌 The Neighbourhood Board

A community bulletin board for neighbours to share announcements, lost & found items, events, and general posts.

**Live:** https://bulletin-board-tuak.onrender.com

## Stack

- **Backend:** Go + Chi router
- **Frontend:** Go HTML templates
- **Database:** SQLite (`modernc.org/sqlite`)
- **Hosting:** Render (with persistent disk)

## Features

- Post to 4 categories: Announcements, Lost & Found, Events, General
- Upvote posts and leave comments
- RSVP to events (Going / Maybe / Not Going)
- Mark Lost & Found items as resolved
- Session-based auth (register, login, logout)
- User profiles

## Project Structure

```
bulletin/
├── cmd/main.go              # Entry point & router
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # DB queries
│   └── middleware/          # Session auth
├── templates/               # Go HTML templates
├── static/css/              # Stylesheet
├── db/schema.sql            # SQLite schema
└── render.yaml              # Render deploy config
```

## Running Locally

```bash
git clone https://github.com/ThatAutocrat/bulletin-board
cd bulletin-board
go mod tidy
go run ./cmd/main.go
```

App runs on http://localhost:8080


### things to add: image uploads, search, map integration, email notifications, admin panel, or pagination. Pull requests will be merged upon approval. Happy coding!
