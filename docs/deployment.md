# Deployment and operations

Pivot's source code is public under the MIT license, while every deployed application remains private and invitation-only by default. Publishing or deploying the code does not authorize third-party data redistribution or Footao collection.

## Production environment

Required values:

- `DATABASE_URL`: pooled PostgreSQL connection string used by the application. TLS is required for hosted databases.
- `MIGRATION_DATABASE_URL`: optional direct PostgreSQL connection string used by Goose. It falls back to `DATABASE_URL` for self-hosted installations.
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

`render.yaml` defines one Render Docker web service named `pivot`. Its multi-stage image builds Vue, bundles the resulting static files with the Go API, and serves both from the same origin. This gives the V1 one public URL (ideally `https://pivot.onrender.com`), avoids third-party session cookies between separate `onrender.com` services, and keeps Vue Router fallbacks inside the application. Render derives `APP_BASE_URL` from `RENDER_EXTERNAL_URL`; set `APP_BASE_URL` explicitly only when attaching a custom domain.

Create a Neon project in a European region, then copy both connection strings from the **Connect** dialog:

- use the pooled hostname containing `-pooler` for `DATABASE_URL`;
- disable pooling and use the direct hostname for `MIGRATION_DATABASE_URL`;
- retain `sslmode=require` and `channel_binding=require` on both URLs.

In the Render Blueprint form, provide:

- `DATABASE_URL`: the Neon pooled connection string.
- `MIGRATION_DATABASE_URL`: the Neon direct connection string.
- `SETUP_TOKEN`: a temporary random value of at least 20 characters known to the operator.
- identifiable user agents and optional provider keys.

The Blueprint requests one free web instance, builds only after GitHub checks pass, and enables remote provider marks for this operator deployment. Free instances can spin down during inactivity, so the first request after an idle period may be slower. Use `SETUP_TOKEN` once at `/setup`, then remove it from the service environment after creating the first administrator. Keep the health check at `/health`. The service binds to Render's `PORT` automatically.

If `pivot.onrender.com` is already allocated, Render assigns another hostname that still contains the service name. A custom `pivot.example.com` domain can be added later without changing the image; set `APP_BASE_URL` to that exact HTTPS origin when doing so.

## Scheduled collections

`.github/workflows/collections.yml` is disabled until the repository variable `COLLECTIONS_ENABLED` equals `true`. Configure the pooled production `DATABASE_URL` and enabled connectors as GitHub Actions secrets, then use `workflow_dispatch` once before relying on schedules. Never store `MIGRATION_DATABASE_URL` in Actions because collection jobs do not perform schema changes.

Default cadence:

- official news hourly;
- matches and standings every six hours;
- Footao daily, only with operator permission;
- club catalog weekly.

The workflow runs server-side Go commands and never contacts providers from a member's browser. GitHub logs contain collection counts and errors but no provider payloads or credentials.

## Release images

Pushing a semantic version tag such as `v1.0.0` runs `.github/workflows/release.yml`. The workflow builds the combined Vue and Go image for AMD64 and ARM64, publishes version and commit tags to `ghcr.io/vincentg2/pivot`, generates an SBOM, and attaches signed build provenance. Remote provider images stay disabled in the distributed image.

```sh
docker pull ghcr.io/vincentg2/pivot:1.0.1
gh attestation verify oci://ghcr.io/vincentg2/pivot:1.0.1 -R vincentg2/pivot
```

The GitHub package may need to be made public once after its first publication. Do not publish `latest` from an unversioned branch build, and never inject runtime credentials as Docker build arguments.
