# Deployment and operations

Pivot remains private by default. Deploying it does not authorize a public repository, third-party data redistribution, or Footao collection.

## Production environment

Required values:

- `DATABASE_URL`: PostgreSQL connection string. TLS is required for hosted databases.
- `APP_BASE_URL`: exact public frontend origin, without a trailing slash.
- `SESSION_SECURE=true`: required behind HTTPS.
- `SETUP_TOKEN`: a random one-time secret of at least 20 characters for the first-install assistant.

Optional connectors use `FOOTBALL_DATA_API_KEY`, `FOOTAO_ENABLED`, `FOOTAO_USER_AGENT`, and `NEWS_USER_AGENT`. Keep all values in the host secret manager. Never use Docker build arguments for credentials.

## Docker Compose

`compose.selfhosted.yml` builds an nginx-served Vue frontend and the Go API, runs PostgreSQL 17, applies migrations on API startup, and exposes the application on port 3000. Replace both required secrets in `.env` before the first start:

```sh
docker compose -f compose.selfhosted.yml up --build -d
docker compose -f compose.selfhosted.yml logs -f api
```

Back up the `pivot-production-data` volume with `pg_dump` before upgrades. Restore into a fresh PostgreSQL instance before starting the new API image. Do not expose PostgreSQL directly to the internet.

## Render and Neon reference

`render.yaml` defines a Docker Go web service and a Render static site. In the Render Blueprint form, provide:

- `DATABASE_URL`: the Neon pooled connection string with `sslmode=require`.
- `APP_BASE_URL`: the final `https://…onrender.com` frontend origin.
- `VITE_API_BASE_URL`: the final API URL ending in `/api/v1`.
- identifiable user agents and optional provider keys.

Render generates `SETUP_TOKEN`; read it from the API service secret environment and use it once at `/setup`. Rotate or remove the value after setup. Keep the API health check at `/health`. The service binds to Render's `PORT` automatically.

## Scheduled collections

`.github/workflows/collections.yml` is disabled until the repository variable `COLLECTIONS_ENABLED` equals `true`. Configure GitHub Actions secrets for the production database and enabled connectors, then use `workflow_dispatch` once before relying on schedules.

Default cadence:

- official news hourly;
- matches and standings every six hours;
- Footao daily, only with operator permission;
- club catalog weekly.

The workflow runs server-side Go commands and never contacts providers from a member's browser. GitHub logs contain collection counts and errors but no provider payloads or credentials.

## Release images

No GHCR publishing workflow is enabled while the repository is private. Add signed, versioned image publication only when the owner explicitly requests the first public release.
