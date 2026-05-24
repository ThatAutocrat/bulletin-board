# 📌 The Neighbourhood Board

A community bulletin board for neighbours to share announcements, lost & found items, events, and general posts.

**Live:** https://bulletin-board-tuak.onrender.com

---

## Stack

- **Backend:** Go + Chi router
- **Frontend:** Go HTML templates
- **Database:** [Turso](https://turso.tech) (libSQL / SQLite-compatible)
- **Hosting:** Render (free tier)

---

## Features

- Post to 4 categories: Announcements, Lost & Found, Events, General
- Upvote posts and leave comments
- RSVP to events (Going / Maybe / Not Going)
- Mark Lost & Found items as resolved
- Session-based auth (register, login, logout)
- User profiles

---

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

---

## Running Locally

```bash
git clone https://github.com/ThatAutocrat/bulletin-board
cd bulletin-board
go mod tidy
go run ./cmd/main.go
```

App runs on http://localhost:8080. A local `bulletin.db` file will be created automatically — no setup needed.

---

## Contributing

Pull requests are welcome! Here are some ideas for things to add:

- Image uploads
- Search
- Map integration
- Email notifications
- Admin panel
- Pagination

PRs will be reviewed and merged upon approval. Happy coding!
