# CipherTalk

A secure, real-time web chat application with direct messaging, group chats, file sharing, and message reactions.

Built as a university project for the [Web and Software Architecture](http://gamificationlab.uniroma1.it/en/wasa/) course at Sapienza University of Rome.

**Live demo:** https://ciphertalk.pages.dev

---

## Features

- **Sign Up / Sign In** — username + password authentication (passwords stored as salted SHA-256 hashes)
- **Direct Messaging** — one-on-one conversations with any registered user
- **Group Chats** — create groups, manage members, assign roles (admin / member)
- **File & Image Sharing** — attach photos and GIFs to messages
- **Message Reactions** — react to messages with a heart
- **Reply & Forward** — thread replies and forward messages to other conversations
- **Read Receipts** — checkmarks show sent / delivered / read status
- **Profile Management** — update username and profile photo
- **User Search** — find other users by username
- **Responsive UI** — desktop sidebar + mobile bottom tab bar

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.21, [httprouter](https://github.com/julienschmidt/httprouter) |
| Database | SQLite (pure Go via [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)) |
| Frontend | Vue 3, Vite, Axios |
| Auth | Bearer token (user ID issued on login) |
| Real-time | Client-side polling every 5 seconds |
| Deployment | Fly.io (backend) + Cloudflare Pages (frontend) |

---

## Project Structure

```
cmd/
  webapi/          # Backend entry point (HTTP server)
demo/
  config.yml       # Local development config
doc/
  api.yaml         # OpenAPI 3.0 specification
service/
  api/             # HTTP handlers and business logic
  database/        # SQLite data access layer
  globaltime/      # Time wrapper (for testing)
vendor/            # Vendored Go dependencies
webui/
  src/
    views/         # Vue page components (Login, Home, Chat, Profile, …)
    components/    # Shared components (Sidebar, etc.)
  public/          # Static assets
Dockerfile.backend
Dockerfile.frontend
fly.toml           # Fly.io deployment config
```

---

## Running Locally

### Prerequisites
- Go 1.21+
- Node.js 18+ and npm

### Backend

```bash
go run ./cmd/webapi/
```

The API server starts on `http://localhost:3000`.

By default it uses the DB path set in `demo/config.yml`. You can override it:

```bash
go run ./cmd/webapi/ --db-filename /tmp/ciphertalk.db
```

### Frontend

```bash
cd webui
npm install
npm run dev
```

The dev server starts on `http://localhost:5173` and proxies API calls to port 3000.

---

## Building for Production

### Backend

```bash
go build -o webapi ./cmd/webapi/
./webapi
```

### Frontend

```bash
cd webui
npm install
npm run build-prod
# Output is in webui/dist/
```

---

## Docker

### Backend

```bash
docker build -f Dockerfile.backend -t ciphertalk-api .
docker run -p 3000:3000 ciphertalk-api
```

### Frontend

```bash
docker build -f Dockerfile.frontend -t ciphertalk-web .
docker run -p 80:80 ciphertalk-web
```

---

## Deployment

### Backend — Fly.io

```bash
fly auth login
fly launch --name ciphertalk-api --dockerfile Dockerfile.backend
fly volumes create db_data --size 1 --region ams
fly secrets set WEBAPI_DB_FILENAME=/data/ciphertalk.db
fly deploy
```

### Frontend — Cloudflare Pages

1. Connect your GitHub repo at [pages.cloudflare.com](https://pages.cloudflare.com)
2. Set **Build command:** `cd webui && npm install && npm run build-prod`
3. Set **Build output directory:** `webui/dist`
4. Add environment variable: `VITE_API_URL` = `https://ciphertalk-api.fly.dev`

---

## API

The full REST API is documented in [`doc/api.yaml`](doc/api.yaml) (OpenAPI 3.0).

Key endpoints:

| Method | Path | Description |
|---|---|---|
| POST | `/session` | Sign in or sign up |
| GET | `/users` | Search users by username |
| GET/PUT | `/users/{id}` | Get or update user profile |
| GET | `/conversations` | List all conversations |
| GET/POST | `/conversations/{id}/messages` | Get or send messages |
| POST | `/groups` | Create a group |
| PUT/DELETE | `/groups/{id}` | Update or delete a group |
| POST/DELETE | `/groups/{id}/members/{uid}` | Add or remove group members |
| POST/DELETE | `/messages/{id}/reactions` | Add or remove a reaction |

---

## Security Notes

- Passwords are hashed with a random 16-byte salt + SHA-256 before storage
- Bearer tokens are currently plain user IDs — not signed JWTs (suitable for academic use)
- No rate limiting is applied — not recommended for high-traffic production use
- HTTPS is enforced by Fly.io and Cloudflare in the deployed environment

---

## Academic Context

This project was developed as part of the **Web and Software Architecture** course (A.Y. 2024–2025) at [Sapienza University of Rome](https://www.uniroma1.it/en/).

The repository structure and backend skeleton are based on the course template by [Prof. Enrico Bassetti](https://github.com/enbas).

---

## License

This project is licensed under the MIT License.

- Original template © 2022 Enrico Bassetti
- CipherTalk implementation © 2024–2025 Aleksandar Kirilov

See [LICENSE](LICENSE) for the full text.
