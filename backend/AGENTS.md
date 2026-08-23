# Backend agent guide

- Organize each domain as handler → service → repository and keep HTTP concerns out of services.
- Use pgx/sqlc-compatible SQL, Goose migrations, `go-playground/validator`, and Echo v4.
- Return RFC 7807 responses through `internal/httpx`; never leak internal errors.
- Authentication uses opaque, hashed database sessions in HttpOnly cookies. Never log tokens or passwords.
- Use `log/slog`, context-aware database calls, explicit timeouts, and graceful shutdown.
- Unit-test services with fakes; PostgreSQL integration tests must be opt-in via `TEST_DATABASE_URL`.
- Keep OpenAPI synchronized with handlers during development.
