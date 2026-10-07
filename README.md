# Ink Mark

A local-first, single-user markdown writing app built with **Go** (backend), **SvelteKit + TypeScript** (frontend), and **PostgreSQL** (everything).

---

## Architecture

``` md
ink-mark/
├── .devcontainer/        # VS Code dev container config
│   ├── devcontainer.json
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── post-create.sh
├── backend/              # Go API server
├── frontend/             # SvelteKit + TypeScript app
└── db/
    └── postgres/
        └── init/         # SQL files — run automatically on first container start
            └── 001_schema.sql
```

## Database

**PostgreSQL only.** Single-user, no sync, no multi-user = no need for a second database.

| Table | Purpose |
|---|---|
| `notebooks` | Top-level organisational groups |
| `documents` | Document title + raw markdown content + `tsvector` for full-text search |
| `tags` / `document_tags` | Tagging system |

Full-text search is handled natively by PostgreSQL via a `tsvector` column and GIN index on `documents`, auto-updated by a trigger on every insert/update.

---

## Getting started

1. Open this repo in VS Code
2. When prompted, click **"Reopen in Container"**
3. The `post-create.sh` script will automatically install Go modules and frontend dependencies

### Services

| Service | URL / Connection |
|---|---|
| Frontend (Vite) | http://localhost:13000 |
| Go API | http://localhost:18080 |
| PostgreSQL | `localhost:15432` — `inkmark` / `inkmark_dev` / db: `inkmark` |
