<p align="center">
  <img src="frontend/public/pivot-mark.svg" width="72" height="72" alt="Pivot logo">
</p>

<h1 align="center">Pivot</h1>

<p align="center">
  A private football dashboard for the clubs you actually follow.
</p>

<p align="center">
  <a href="https://github.com/vincentg2/pivot/actions/workflows/ci.yml"><img src="https://github.com/vincentg2/pivot/actions/workflows/ci.yml/badge.svg" alt="CI status"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-2457ff.svg" alt="MIT license"></a>
</p>

Pivot brings fixtures, results, standings, French TV listings, and official club news into one responsive, invitation-only experience. It is designed for a small private community first and packaged as a self-hostable MIT-licensed project.

> **Résumé français —** Pivot est un tableau de bord football privé et responsive. Il rassemble les clubs favoris, matchs, résultats, classements, diffusions françaises et actualités officielles. L'inscription se fait sur invitation et le projet peut être auto-hébergé.

## What Pivot includes

- An invitation-only account system with opaque database sessions, secure cookies, roles, password recovery, and account deletion.
- Personal profiles with synchronized light, dark, and system themes, plus up to five favorite clubs.
- A month-ahead dashboard with favorite-club filters, recent results, team marks, and TV channels when announced.
- Club pages, fixtures, results, national standings, and date-based match navigation.
- French broadcast listings with channel marks, cautious fixture matching, admin corrections, and an audit trail.
- Configurable official RSS and Atom feeds that retain headline metadata only.
- French by default, with a complete English interface.
- Administration for invitations, collectors, news sources, TV corrections, and one-time password-reset links.

Pivot remains useful without any third-party key: the application starts with an empty catalog and every connector is an explicit operator choice.

## Stack

```text
frontend/   Vue 3 · TypeScript · Vite · Tailwind CSS 4 · Reka UI · Pinia
backend/    Go · Echo 4 · pgx · Goose · PostgreSQL · structured slog
platform/   Docker Compose locally · Render + Neon reference deployment
```

The backend is organized by domain using `handler → service → repository`. The frontend keeps shared state deliberately small: Pinia is limited to session and favorites. API failures use RFC 7807 problem details, and sessions are opaque, hashed, and stored in PostgreSQL.

## Quick start

### Requirements

- Go 1.24+
- Node.js 20+
- pnpm
- Just
- Docker with Compose

### Local development

```sh
cp .env.example .env
just doctor
just install
just db-up
just migrate
just seed-demo
just dev
```

Open [http://localhost:5173](http://localhost:5173). Vite proxies `/api` to the Go API at `http://localhost:8080`; PostgreSQL is the only service running in Docker. Demo seeding prints a one-time invitation code.

Useful commands:

```sh
just check                              # lint, unit tests, builds
just create-admin email@example.com me # create an administrator
just collect-catalog                    # clubs and competitions
just collect-sport                      # matches and standings
just collect-tv                         # French TV listings
just collect-news                       # official club headlines
```

## First installation

On an empty database, Pivot redirects to `/setup`.

1. Generate a random `SETUP_TOKEN` of at least 20 characters.
2. Add it to the runtime environment before starting Pivot.
3. Open `/setup`, enter the token, and create the first administrator.
4. Remove `SETUP_TOKEN` from the runtime environment after setup.

The setup route locks permanently once a user exists. `just create-admin` remains available as an operator recovery path.

## Data connectors

| Source | Provides | Default | Operator responsibility |
| --- | --- | --- | --- |
| [football-data.org](https://www.football-data.org/) | Clubs, fixtures, results, standings | Disabled without a key | Supply `FOOTBALL_DATA_API_KEY` and keep attribution visible |
| Footao | French TV listings | Disabled | Enable only with permission and an identifiable User-Agent |
| Official club RSS/Atom | Headline, source, date, link | No feeds configured | Add public official feeds in the admin interface |
| Remote provider images | Club and channel marks | Disabled | Set `VITE_REMOTE_LOGOS_ENABLED=true`; images are never distributed with Pivot |

Footao is fetched centrally by the server—never by a member's browser—with bounded retries and a two-month window. Past TV data is removed the following day. Official-news collection blocks private-network targets and keeps metadata for 30 days.

See [Data sources and retention](docs/data-sources.md) and [Third-party notices](THIRD_PARTY_NOTICES.md) before enabling any connector. The MIT license covers Pivot's code only, not third-party data, names, or images.

## Self-hosting

For a complete local installation, copy `.env.example`, replace `POSTGRES_PASSWORD` and `SETUP_TOKEN`, then run:

```sh
docker compose -f compose.selfhosted.yml up --build -d
docker compose -f compose.selfhosted.yml logs -f api
```

Open [http://localhost:3000/setup](http://localhost:3000/setup). Use TLS and `SESSION_SECURE=true` outside localhost. The API container applies Goose migrations before it starts.

The reference hosted setup uses one Render web service for Vue and Go, plus Neon PostgreSQL. Scheduled GitHub Actions collect news hourly, sports data every six hours, TV listings daily, and the club catalog weekly when explicitly enabled by the operator.

Read [Deployment and operations](docs/deployment.md) for production secrets, pooled and direct database URLs, backups, health checks, and scheduled collectors.

## Repository guide

- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Deployment and operations](docs/deployment.md)
- [Data sources and retention](docs/data-sources.md)
- [Third-party notices](THIRD_PARTY_NOTICES.md)
- [MIT license](LICENSE)

## Project status

The five initial product milestones are implemented. The repository stays private until its owner explicitly requests the first public release; GHCR image publishing will begin with that release.
