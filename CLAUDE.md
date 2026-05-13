# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

### Backend (Gin + MongoDB)
`server/main.go` wires everything: CORS, cookie sessions, Mongo init (`db.Init`), Google Cloud Tasks init (`services/gcloud.InitTasks`), then mounts API groups under `/api` via `routes.Init*` and `slackbot.InitSlackbot`. After API routes, it walks `frontend/dist` and registers each file as a static route, loads `index.html` as a template, and falls back to a `NoRoute` handler that injects per-route OG meta tags (e.g. for `/e/:eventId` it looks up the event to set the title and OG image).

- `routes/` — HTTP handlers grouped by domain: `auth.go`, `user.go`, `users.go`, `events.go`, `folders.go`, `analytics.go`, `stripe.go`. Route comments use Swag annotations; `swag init` regenerates `docs/`.
- `models/` — Mongo document structs (`Event`, `User`, `Response`, `Folder`, `Attendee`, `Calendar`, `Set`, `Otp`, `FriendRequest`, `Location`, `DailyUserLog`).
- `db/` — Mongo accessors per model (`events.go`, `users.go`, `folders.go`, `analytics.go`, `utils.go`) plus `init.go`. Treat this as the only layer that talks to Mongo.
- `services/` — external integrations. Notable: `calendar/` (Google, Outlook/Graph, Apple CalDAV via `jonyTF/go-webdav`, generic ICS), `auth/`, `contacts/`, `gcloud/` (Cloud Tasks for scheduled jobs), `listmonk/`, `microsoftgraph/`.
- `middleware/auth.go` — session-based auth middleware applied selectively by `routes.Init*`.
- `slackbot/` and `discord_bot/` — bot integrations registered as additional handlers.
- `scripts/` — one-off Mongo migrations (dated folders like `20250417_responses_collection`). Run manually; don't import from runtime code.
- `utils/` — generic helpers (`array_utils`, `db_utils`, `mail_utils`, `request_utils`, `response_utils`).
- `logger/` — wraps log file (`logs.log`) + stdout via `gin.DefaultWriter`.

### Frontend (Vue 2 SPA)
- `src/router/index.js` — routes (`Landing`, `Home`, `Event`, `Group`, `Friends`, `Settings`, `SignIn`/`SignUp`/`Auth`, `StripeRedirect`, etc. — see `src/views/`).
- `src/store/index.js` — single Vuex store (auth user, events, snackbar, dialogs).
- `src/components/` — organized by feature folder (`event/`, `groups/`, `home/`, `landing/`, `pricing/`, `settings/`, `schedule_overlap/`, `calendar_permission_dialogs/`, `sign_up_form/`, `general/`) plus top-level shared components.
- `src/utils/` — date math (`date_utils.js`, uses `dayjs`/`moment`/`spacetime`), `fetch_utils.js` (API client), `plugin_utils.js` (handles the postMessage plugin API — see `PLUGIN_API_README.md`), `sign_in_utils.js`, `location_utils.js`, `services/` (calendar-provider abstractions on the client side).

## Conventions worth knowing

- The Go module path is `schej.it/server`; imports use that prefix throughout. Don't rename.
- Mongo collection naming and indexes are established by the dated migration scripts in `server/scripts/` — when adding a new collection or index, follow the same dated-folder pattern.
- New API routes need Swag comments above the handler so `swag init` picks them up; otherwise they're invisible in `/swagger`.
- The server panics on startup if `SESSION_SECRET` is missing or shorter than 32 chars (`validateSessionSecret` in `main.go`).
- `frontend/dist` is consumed by the Go server at runtime — local server boot tries `./frontend/dist` then `../frontend/dist`, or honors `FRONTEND_DIST` env var.
