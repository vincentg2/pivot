# Pivot

Pivot is a private, invitation-only football dashboard for friends and a self-hostable open-source portfolio project. It combines favorite clubs, fixtures, results, standings, French TV listings, and official club news in a calm editorial interface.

> **Français :** Pivot est un tableau de bord football privé et responsive, accessible sur invitation. Le projet est auto-hébergeable ; les connecteurs de données restent optionnels et chaque opérateur fournit ses propres accès.

## Milestones 1–5

The foundation includes the Vue/Go monorepo, responsive theme system, PostgreSQL, opaque cookie sessions, invitation registration, profiles, account deletion, administration, tests, CI, and development tooling. The product includes a local club catalog, five synchronized favorites, fixtures, results, standings, French TV listings, official club headlines, and operator-controlled collections. TV-only matches remain visibly external, while cautious matching enriches recognized catalog fixtures.

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

Administrators can issue a one-time password reset link from the administration page. Links expire after 30 minutes, are stored only as hashes, and revoke all existing sessions when consumed.

The application works with an empty catalog and no provider key. To populate the five major European leagues, set `FOOTBALL_DATA_API_KEY` in `.env`, restart the API, then use **Admin → Data collection → Run now**. Remote crest rendering remains off unless the operator explicitly sets `VITE_REMOTE_LOGOS_ENABLED=true`; URLs are stored, but images are never distributed with Pivot.

After the club catalog is populated, run **Admin → Sports collection → Run sports data** or `just collect-sport`. The importer deliberately paces provider requests, stores yesterday through the next 30 days, refreshes current standings, and preserves previously encountered seasons and matches.

The Footao connector is a separate, explicit opt-in. Only enable it when the operator has permission: set `FOOTAO_ENABLED=true` and replace the placeholder `FOOTAO_USER_AGENT` with an identifiable contact URL, restart the API, migrate, then use **Admin → Television collection → Run TV data** or `just collect-tv`. Pivot performs one central server-side fetch with bounded retries, imports a two-month schedule window, and never contacts Footao from browsers. Manual listing corrections and restores are recorded in the admin audit trail.

Official club news is configured per club under **Admin → Official news**. Pivot accepts public RSS and Atom URLs, blocks private-network fetches, and stores only titles, source names, publication dates, and canonical links for 30 days. Run `just collect-news` or enable the opt-in scheduled workflow described in [deployment documentation](docs/deployment.md).

## First installation

On an empty database, Pivot redirects to `/setup`. Set a secret `SETUP_TOKEN` of at least 20 characters before starting the API, enter it once in the web assistant, and create the first administrator. The route locks permanently as soon as any user exists. Existing installations are unaffected, and `just create-admin` remains available for operator recovery.

## Self-hosting with Docker Compose

The default `docker-compose.yml` still runs PostgreSQL only for local development. For a complete installation, copy `.env.example` to `.env`, replace `POSTGRES_PASSWORD` and `SETUP_TOKEN`, then run:

```sh
docker compose -f compose.selfhosted.yml up --build -d
```

Open `http://localhost:3000/setup`. Use TLS and `SESSION_SECURE=true` outside localhost. The API container applies Goose migrations before starting. See [deployment documentation](docs/deployment.md) for backups, Neon, Render, and scheduled collections.

## Data and licensing

Pivot works without data-provider keys. Read [data source policy](docs/data-sources.md) and [third-party notices](THIRD_PARTY_NOTICES.md) before enabling connectors. The code is MIT licensed; third-party data, names, and images are not.

## Status

Milestone 5 of 5. The repository remains private until the owner explicitly decides to publish it. GHCR publication intentionally starts only when the owner requests the first public release.
