# Contributing to Pivot

Pivot is in early development. Before opening a pull request, discuss substantial product or schema changes in an issue.

1. Install Go 1.24+, Node 20+, pnpm 10+, Just, and Docker with Compose.
2. Copy `.env.example` to `.env`, then run `just doctor` and `just db-up`.
3. Run the API and frontend with `just dev`.
4. Before submitting, run `just check`.

Use conventional, imperative commit subjects. Add migrations for schema changes, tests for behavior, and documentation for operational changes. Contributions must not contain third-party data or secrets and are accepted under the MIT license.
