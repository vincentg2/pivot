# Security policy

## Reporting

Do not open a public issue for a vulnerability. Send a private GitHub security advisory to the maintainers with reproduction steps, impact, and any proposed remediation. Expect an acknowledgement within seven days.

## Supported versions

Until the first stable release, only the latest commit on `main` receives security updates.

## Deployment baseline

Use TLS, set `SESSION_SECURE=true`, rotate all sample credentials, restrict database access, and keep dependencies patched. Never expose `.env`, invitation codes, reset links, session tokens, or collection credentials in logs.
