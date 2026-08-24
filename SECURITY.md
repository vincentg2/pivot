# Security policy

## Reporting

Do not open a public issue for a vulnerability. Send a private GitHub security advisory to the maintainers with reproduction steps, impact, and any proposed remediation. Expect an acknowledgement within seven days.

## Supported versions

Until the first stable release, only the latest commit on `main` receives security updates.

## First-install safety

Always set a unique `SETUP_TOKEN` before exposing an empty installation. Complete `/setup` immediately, verify the first administrator, then rotate or remove the token. Creation is transactionally limited to an empty `users` table, but the token still prevents an unauthenticated race during initial deployment.

## Deployment baseline

Use TLS, set `SESSION_SECURE=true`, rotate all sample credentials, restrict database access, and keep dependencies patched. Never expose `.env`, invitation codes, reset links, session tokens, or collection credentials in logs.
