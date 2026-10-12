# Database Migrations

Managed with [Goose](https://pressly.github.io/goose/) v3. The database is **SQLite** — the `.db` file lives next to the binary at runtime (`ink-mark.db`).

Migrations run automatically at app startup. For manual use during development:

## Commands

```sh
# Run all pending migrations
goose -dir db/migrations sqlite3 ./ink-mark.db up

# Roll back the last migration
goose -dir db/migrations sqlite3 ./ink-mark.db down

# Check current migration status
goose -dir db/migrations sqlite3 ./ink-mark.db status

# Create a new migration file
goose -dir db/migrations create <name> sql
```

## File naming

Goose uses sequential numeric prefixes:

```
00001_create_documents.sql
00002_add_tags.sql
```

Run all commands from the **project root** (`/workspace/ink-mark`).
