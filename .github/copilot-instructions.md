# Copilot Instructions for SoditApp

This file helps Copilot-style assistants quickly understand how to build, run, and contribute to this repository.

## Quick build / dev commands
- Install frontend deps: `cd frontend && pnpm install`
- Live development (frontend hot reload + Go backend): `wails dev`
- Build production distributable: `wails build` (binary in `build/bin/`)

Frontend scripts (from frontend/package.json):
- `pnpm run dev` — vite dev server
- `pnpm run build` — production frontend build
- `pnpm run preview` — preview built frontend
- `pnpm run check` — `svelte-check` + `tsc` (type + Svelte checks)

Go / tests / linting:
- Run all Go tests: `go test ./...`
- Run tests for a single package: `go test ./path/to/package -v`
- Run a single test by name: `go test ./... -run '^TestName$' -v`
- Lightweight vet: `go vet ./...`

Database / migrations:
- SQL migrations and seed files live under `db/migrations/` and `db/seeds/`.

## High-level architecture (short)
- Desktop app using Wails v2: Go backend + embedded web frontend (Svelte 5 + Vite).
- Data flow: Frontend → Wails-generated TypeScript bindings (frontend/wailsjs) → Go bound methods → service layer → repository layer → PostgreSQL.
- Project layout (important folders):
  - `main.go`, `app.go` — Wails bootstrap and app lifecycle
  - `backend/` — (intended) domain logic: models, repository, services, dto
  - `frontend/` — Svelte 5 app (src/, lib/ components, shadcn-svelte primitives)
  - `db/` — migrations and seed SQL files
  - `docs/` — requirements & design (see `docs/soditapp.md`)

## Key conventions for this repo
- UI text is Spanish. Code remains English.
- Frontend uses Svelte 5 runes only: prefer `$state()`, `$derived()`, `$effect()`, `$props()`, `$bindable()`; do NOT use legacy `$:` or Svelte 4 stores.
- All `.svelte` files must use `<script lang="ts">` and explicit types; use `frontend/src/lib/types/` for shared types.
- Use `pnpm` for frontend package management (configured in `wails.json`).
- Use `shadcn-svelte` primitives in `src/lib/components/ui/`; these are managed by the shadcn tool—avoid manual edits unless intentional.
- Wails bindings: only expose service-level methods (on `App` or structs registered in `Bind`). Do not bind repository-layer internals directly.
- DTOs (DeliverySheetDTO, DeliveryEntryDTO) are ephemeral and must NOT be persisted; they exist solely for print/preview.
- Balance mutation rule: only update `ClientProduct.balance` via delivery confirmation logic (`balance += credit_units`); recalc deltas on edits/deletes.
- Delivery sheets: exactly 4 day columns, grouped by zone, only show clients with negative balances; paginate at A4 boundaries.

## Where to read more
- `docs/soditapp.md` — requirements, entities, balance logic, print layout.
- `AGENTS.md` — project-specific coding & AI agent guidance (Go + Svelte rules).
- `CLAUDE.md` — additional assistant-oriented build and convention notes.

## Other assistant configs discovered
- AGENTS.md and CLAUDE.md are present and contain canonical rules for Go & Svelte 5 in this repo.
- No `.cursorrules`, `.windsurfrules`, `AIDER_CONVENTIONS.md`, or `.clinerules` were found at time of creation.

---

If anything here should be expanded (example: exact go package paths for tests, custom linters, or CI workflow commands), say which area and more details will be added.
