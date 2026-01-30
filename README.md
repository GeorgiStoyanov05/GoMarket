# GoMarket

GoMarket is a Go web application structured in an MVC-ish style (controllers / routes / services / models) with server-rendered views and static assets.

---

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration (Environment Variables)](#configuration-environment-variables)
  - [Run](#run)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- Server-rendered UI (`views/`) with static assets (`static/`)
- Organized backend structure: controllers, routes, services, models
- Database layer (`database/`)
- Middleware layer (`middlewares/`)
- Common app flows (examples):
  - authentication + session/JWT handling
  - user actions (e.g., deposits / portfolio actions)
  - market data pages (quotes, charts, watchlists)

---

## Project Structure

```
.
├── controllers/   # Request handlers / controller logic
├── routes/        # Route definitions + wiring
├── services/      # Business logic / integrations
├── models/        # Domain models / DTOs
├── middlewares/   # Auth/logging/etc.
├── database/      # DB connection + queries/repositories
├── views/         # HTML templates (server-side)
├── static/        # CSS/JS/images
├── main.go        # App entrypoint
├── go.mod
└── go.sum
```

---

## Tech Stack
- **Backend:** Go
- **Frontend:** HTML + CSS + JS (server-rendered templates + static assets)
- **Architecture:** controllers / services / routes separation
- **Database:** configured via `database/` (check your implementation for the engine)

---

## Getting Started

### Prerequisites
- Go installed (recommended: latest stable)
- A database instance if your app requires one (see `database/`)
- Any API keys required by your services (if applicable)

### Installation
```bash
git clone https://github.com/GeorgiStoyanov05/GoMarket.git
cd GoMarket
go mod download
```

### Configuration (Environment Variables)
Create a `.env` file (or export env vars in your shell).  
**Rename these to match what your code expects** (look in `main.go`, `database/`, and `services/`):

```env
# Server
PORT=8080

# Database (examples)
DB_URL=postgres://user:pass@localhost:5432/gomarket?sslmode=disable
# or
MONGO_URI=mongodb://localhost:27017/gomarket

# Auth (examples)
JWT_SECRET=replace_me
SESSION_SECRET=replace_me

# External services (examples)
MARKET_DATA_API_KEY=replace_me
```

> Tip: If you don’t use `.env`, remove that part and just export variables normally.

### Run
```bash
go run .
```

Then open:
- `http://localhost:8080` (or whatever `PORT` is set to)

---

## Development

Common commands:
```bash
go fmt ./...
go test ./...
go vet ./...
```

Build:
```bash
go build -o gomarket
./gomarket
```

---

## Troubleshooting

**1) “missing env var” / app crashes at startup**  
- Check which variables your code reads in `main.go`, `database/`, and `services/`.
- Add them to `.env` or export them before running.

**2) Templates not loading / blank pages**  
- Verify `views/` paths are correct.
- Ensure your template loader uses the right working directory (running from repo root usually helps).

**3) Static files not loading (CSS/JS missing)**  
- Confirm your router serves `static/` (e.g. `/static/...`).
- Check the `<link>` / `<script>` paths in your templates.

**4) Port already in use**  
```bash
export PORT=8090
go run .
```

---

## License
No license is currently specified.
