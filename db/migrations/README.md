# Database Migrations

Managed with [Goose](https://pressly.github.io/goose/) v3.

## Connection

```sh
export DATABASE_URL=postgres://inkmark:inkmark_dev@postgres:5432/inkmark?sslmode=disable
```

This is already set in the devcontainer environment.

## Commands

```sh
# Run all pending migrations
goose -dir db/migrations postgres "$DATABASE_URL" up

# Roll back the last migration
goose -dir db/migrations postgres "$DATABASE_URL" down

# Check current migration status
goose -dir db/migrations postgres "$DATABASE_URL" status

# Create a new migration file
goose -dir db/migrations create <name> sql
```

## File naming

Goose uses sequential numeric prefixes:

```
00001_create_documents.sql
00002_add_notebooks.sql
```

Run all commands from the **project root** (`/workspace/ink-mark`).
