# Data sources and retention

## football-data.org

The primary catalog connector supports competitions available to the operator's plan. Each installation supplies its own API key. With no key, Pivot stays usable for authentication, profiles, manually configured favorites, and administration. Provider attribution must remain visible wherever its data appears.

Remote badge display is an operator setting. Images are linked from the provider and never copied into the repository. A generated monogram is the default fallback.

Suggested schedule: club catalog weekly; fixtures and standings four times daily. Past seasons are retained.

## Footao

The French TV listing connector is centrally fetched by the backend, disabled by default, and must never be called from a browser. Enabling it confirms the operator has permission. Requests use an identifiable User-Agent, conservative backoff, and stop on repeated failures. Audited admin corrections override imported values.

Suggested schedule: daily. Past TV listings are deleted the following day. TV-only matches outside the football catalog remain explicitly external and unenriched.

## Official club RSS

Feeds are operator-configurable. Pivot stores only title, source, date, and link for 30 days, with an hourly suggested schedule.

## Match windows

The all-matches page is designed around Yesterday (results only, no historical TV), Today, Tomorrow, the next seven days, arbitrary date navigation, and filters.
