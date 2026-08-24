# Data sources and retention

## football-data.org

The primary catalog connector supports competitions available to the operator's plan. Each installation supplies its own API key. With no key, Pivot stays usable for authentication, profiles, manually configured favorites, and administration. Provider attribution must remain visible wherever its data appears.

Remote badge display is an operator setting (`VITE_REMOTE_LOGOS_ENABLED=false` by default). Images are linked from the provider and never copied into the repository. A generated monogram is the default fallback. Catalog collection is launched from the administration screen and its latest outcome is retained locally.

Suggested schedule: club catalog weekly; fixtures and standings four times daily. Past seasons are retained.

The sports collector fetches yesterday through the next seven days and the current total standings for the five major leagues. Requests are paced for conservative free-plan usage. Repeated runs upsert matches and tables while leaving older seasons intact.

## Footao

The French TV listing connector is centrally fetched by the backend, disabled by default, and must never be called from a browser. Enabling it confirms the operator has permission. Requests use an identifiable User-Agent, conservative backoff, and stop on repeated failures. Audited admin corrections override imported values.

Suggested schedule: daily. Past TV listings are deleted the following day. TV-only matches outside the football catalog remain explicitly external and unenriched.

Activation requires both `FOOTAO_ENABLED=true` and an operator-specific `FOOTAO_USER_AGENT`. The checked-in placeholder is intentionally rejected when the connector is enabled. A collection imports today through the next two calendar months from one response, stores normalized listing fields rather than raw HTML, and discards the response after parsing. Listings are linked to football-data.org matches only when kickoff time and both team names provide an unambiguous match. Administrators can correct or hide a listing; the effective correction survives later imports and every correction or restore creates an immutable audit entry.

## Official club RSS

Feeds are operator-configurable. Pivot stores only title, source, date, and link for 30 days, with an hourly suggested schedule.

Only public HTTP and HTTPS endpoints are accepted. The server rejects localhost, private, link-local, multicast, and credential-bearing URLs before fetching and revalidates redirects. Feed responses are limited to 2 MiB. Article bodies, summaries, media, authors, and raw XML are discarded.

## Match windows

The all-matches page is designed around Yesterday (results only, no historical TV), Today, Tomorrow, the next seven days, arbitrary date navigation, and filters.
