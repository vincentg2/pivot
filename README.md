# Pivot

Pivot is a private, invitation-only football dashboard for friends and a self-hostable open-source portfolio project. It combines favorite clubs, fixtures, results, standings, French TV listings, and official club news in a calm editorial interface.

> **Français :** Pivot est un tableau de bord football privé et responsive, accessible sur invitation. Le projet est auto-hébergeable ; les connecteurs de données restent optionnels et chaque opérateur fournit ses propres accès.

## Milestones 1–4

The foundation includes the Vue/Go monorepo, responsive theme system, PostgreSQL, opaque cookie sessions, invitation registration, profiles, account deletion, administration, tests, CI, and development tooling. The product now includes a local club catalog, club profiles, five synchronized favorites, fixtures and results by date, current domestic standings, French TV listings, and operator-controlled data collections. TV-only matches remain visibly external, while cautious matching enriches recognized catalog fixtures.

## Architecture

- `frontend/`: Vue 3, TypeScript, Vite, Tailwind CSS v4, Reka UI, Lucide Vue, and Pinia for session/favorites only.
- `backend/`: Go, Echo v4, pgx, Goose migrations, validator, structured logging, and domain-oriented handler/service/repository packages.
- PostgreSQL is the only stateful dependency. Go and Vue run locally during development.

## Developer setup

Requirements: Go 1.24+, Node 20+, pnpm, Just, and Docker Compose.

```sh
cp .env.example .env
just doctor
just install
just db-up
just migrate
just seed-demo
just dev
```

The frontend runs at `http://localhost:5173` and proxies `/api` to the backend at `http://localhost:8080`. Demo seeding prints a one-time invitation code. Production credentials are never generated or committed.

Run all checks with `just check`. Create an admin interactively with `just create-admin email@example.com nickname`.

The application works with an empty catalog and no provider key. To populate the five major European leagues, set `FOOTBALL_DATA_API_KEY` in `.env`, restart the API, then use **Admin → Data collection → Run now**. Remote crest rendering remains off unless the operator explicitly sets `VITE_REMOTE_LOGOS_ENABLED=true`; URLs are stored, but images are never distributed with Pivot.

After the club catalog is populated, run **Admin → Sports collection → Run sports data** or `just collect-sport`. The importer deliberately paces provider requests, stores yesterday through the next seven days, refreshes current standings, and preserves previously encountered seasons and matches.

The Footao connector is a separate, explicit opt-in. Only enable it when the operator has permission: set `FOOTAO_ENABLED=true` and replace the placeholder `FOOTAO_USER_AGENT` with an identifiable contact URL, restart the API, migrate, then use **Admin → Television collection → Run TV data** or `just collect-tv`. Pivot performs one central server-side fetch with bounded retries; browsers never contact Footao. Manual listing corrections and restores are recorded in the admin audit trail.

## Self-hosting with Docker Compose

The checked-in Compose file intentionally runs PostgreSQL only for local development. A production Compose profile and first-install web assistant will land before the first public release. Until then, build `frontend/` as static assets, run the Go API behind TLS, set a strong PostgreSQL password, and set `SESSION_SECURE=true`.

## Data and licensing

Pivot works without data-provider keys. Read [data source policy](docs/data-sources.md) and [third-party notices](THIRD_PARTY_NOTICES.md) before enabling connectors. The code is MIT licensed; third-party data, names, and images are not.

## Status

Milestone 4 of 5. The repository remains private until the owner explicitly decides to publish it.
