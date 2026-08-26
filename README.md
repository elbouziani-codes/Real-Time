# Real-Time Forum

A real-time forum with posts, comments, reactions, categories and private chat, built with a Go (SQLite) backend and a vanilla JavaScript SPA frontend.

## Run

```sh
CGO_ENABLED=1 go run cmd/main.go
```

- App: http://localhost:8081
- Run from the repo root (expects `./web` and `./schema.sql`)
- DB (`realTime.db`) is auto-created and migrated at startup; only the 4 default categories are seeded — register a user via `/register`

## Stack

- **Backend**: Go 1.27, `gorilla/websocket`, `mattn/go-sqlite3`, `golang.org/x/crypto`
- **Frontend**: vanilla JS SPA (history-API routing), served statically by the backend with an `index.html` fallback

## Architecture

```
cmd/main.go                 entry point & manual dependency wiring
internal/
  config/                   port (:8081) + DB path
  domain/                   entities, validation, error codes
  repository/               SQL queries per entity (sqlite driver in sqlite/)
  service/                  business logic
  transport/http/           handlers, auth middleware, router
  transport/webSocket/      chat hub & presence
web/                        SPA (index.html, css/, js/)
```

## API

All routes except register/login require session auth.

| Method | Route | Purpose |
|---|---|---|
| POST | `/api/register` | create account |
| POST | `/api/login` | login (email or nickname) |
| POST | `/api/logout` | logout |
| GET | `/api/me` | current user |
| GET/POST | `/api/posts` | list / create posts |
| GET | `/api/posts/{id}` | post details |
| POST | `/api/comments` | create comment |
| GET | `/api/comments/{id}` | comments of a post |
| GET/POST/PATCH | `/api/reactions...` | like/dislike posts & comments |
| GET | `/api/categories` | list categories |
| GET | `/api/users`, `/api/users/{id}` | users (paged) / profile |
| POST | `/api/getMessage` | chat history paging |
| GET | `/api/ws` | WebSocket |

## WebSocket protocol (`/api/ws`)

Client → server: `{ request_type, mod, id, content, destination }` where `request_type` is `"message"` or `"typing"`.

Server → client codes:

| Code | Meaning |
|---|---|
| 1 | connection replaced by another tab |
| 2 | typing indicator |
| 3 | user came online (new clients also get the online list) |
| 4 | user went offline |
| 200 | chat message delivered |
| 404 | request failed |

## Frontend structure (`web/js`)

- `router/` — history-API SPA router (`/`, `/postDetails`, `/chat`, `/login`, `/register`, `/end`)
- `pages/`, `components/`, `listeners/` — render functions, UI pieces, per-page event wiring
- `services/` — state layer (users Map as single source of truth for presence/conversations)
- `api/` — fetch wrappers
- `websocket/socket.js` — singleton WS connection; events routed through `services/messages.js`
