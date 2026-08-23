# Pivot agent guide

## Scope

These rules apply to the whole monorepo. Read the nearest nested `AGENTS.md` before changing frontend or backend code.

## Product invariants

- Pivot is private by default. Never weaken invitation-only registration or expose secrets.
- The application must remain useful without third-party API keys.
- Do not commit third-party data, club badges, API responses, credentials, or personal data.
- Treat football-data.org attribution and connector opt-ins as product requirements.
- Keep the interface accessible, responsive down to 375 px, and usable in system, light, and dark themes.

## Engineering workflow

- Use `pnpm` for frontend dependencies and Go modules for backend dependencies.
- Prefer focused domain packages and explicit dependencies over global state.
- Add or update tests with behavior changes. Run `just check` before committing.
- Database changes require a Goose migration and corresponding query/repository changes.
- API errors use RFC 7807 problem details. Logs are structured and must not contain secrets.
- Keep commits small, descriptive, and free of generated local state.

## Repository hygiene

- `.env.example` contains safe placeholders only. Real `.env` files are ignored.
- `CLAUDE.md` files must remain relative symlinks to their neighboring `AGENTS.md`.
- English is the canonical documentation language; the README includes a short French summary.
