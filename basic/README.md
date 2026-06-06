# basic

A small, production-structured Go backend used as a learning reference. Two
domains — `user` and `post` (a post belongs to a user) — each wired through the
standard layers:

```
HTTP (Gin handler)  ->  Service (business rules)  ->  Repository (storage)  ->  Postgres
```

## Layout

```
cmd/server/main.go            entry point: load config, open db, wire layers, run
cmd/migrate/main.go           schema migrations: up / down / version
internal/
  platform/config            env -> typed Config
  platform/database          opens the GORM connection; runs migrations
  platform/database/migrations  versioned .sql files (embedded), the source of truth for schema
  user/                       model, repository (interface + gorm impl), service, handler
  post/                       same, plus a cross-domain check that the author exists
  server/router.go            builds the gin.Engine and mounts the routes
```

Key ideas to look at:

- **`internal/`** is private to this module (enforced by the go tool).
- **Dependency injection** lives in `cmd/server/main.go` — nothing uses a global DB.
- **Repositories are interfaces**, so the service is testable without a database
  (`internal/post/service_test.go`).
- **`post.UserChecker`** is a consumer-side interface: posts depend on a tiny
  slice of the user service, not the whole package.

## Run

```sh
cp .env.example .env          # point DATABASE_URL at a local Postgres
go run ./cmd/migrate up       # create/upgrade the schema (run once, and after each deploy)
go run ./cmd/server           # start the API
```

The schema is owned by the SQL files in `internal/platform/database/migrations`,
not by the app. The server never alters tables on startup — migrations are a
deliberate, separate, reversible step.

```sh
go run ./cmd/migrate up        # apply all pending migrations
go run ./cmd/migrate down      # roll back the most recent migration
go run ./cmd/migrate version   # show current schema version
```

To add a schema change, create the next numbered pair, e.g.
`0003_add_posts_published.up.sql` and `0003_add_posts_published.down.sql`.

## Test

```sh
go test ./...
```

## Endpoints

```
GET  /health
POST /api/v1/users            {"name","email"}
GET  /api/v1/users
GET  /api/v1/users/:id
POST /api/v1/posts            {"user_id","title","body"}
GET  /api/v1/posts            optional ?user_id=N
GET  /api/v1/posts/:id
```
