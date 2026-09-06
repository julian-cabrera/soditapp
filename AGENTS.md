# AGENTS.md — SoditApp

Desktop delivery/client management app. **Wails v2** (Go backend + Svelte 5 frontend), Windows 10+, PostgreSQL.

Full requirements in `docs/soditapp.md`. Claude-specific guidance in `CLAUDE.md`.

---

## Build & Dev Commands

```bash
wails dev                  # full app: Go backend + Svelte hot reload
wails build                # production binary → build/bin/soditapp

cd frontend && pnpm install   # install frontend deps
pnpm run dev                  # vite only (no Go)
pnpm run build                # vite production build
pnpm run check                # svelte-check + tsc (run before committing frontend changes)

pnpm dlx shadcn-svelte@latest add <component>   # add shadcn component
```

No Go test framework is configured. No CI workflows exist.

---

## Architecture

```
main.go          — Wails bootstrap. Embeds frontend/dist, db/migrations/*.sql, db/seeds/*.sql.
app.go           — App struct. Allocates services at construction, calls Init() in startup().
backend/
  model/         — Domain structs (Zone, Client, Product, ClientProduct, DeliveryTransaction). Singular package name.
  repository/    — SQL access layer. One repo per entity. Constructor: NewXRepository(db).
  service/       — Business logic. Bound to Wails via main.go Bind array.
  dto/           — DeliverySheetDTO, DeliveryEntryDTO. Runtime-only, never persisted.
  database/      — Connect(), RunMigrations(), RunSeeds(). Custom migration tracker (schema_migrations table).
db/migrations/   — Numbered .sql files (001_*.sql, 002_*.sql, ...). Embedded via //go:embed.
db/seeds/        — Numbered .sql files. Embedded via //go:embed. Same tracking pattern as migrations.
frontend/
  src/components/         — App-level pages and dialogs, organized by entity (zones/, clients/, products/, delivery/)
  src/lib/components/ui/  — shadcn-svelte primitives (managed by CLI, do not edit manually)
  src/lib/stores/         — Global state as .svelte.ts files using Svelte 5 runes ($state)
  src/lib/hooks/          — Reusable Svelte 5 hooks
  src/lib/types/          — Shared TypeScript types
  src/lib/utils.ts        — cn() for Tailwind class merging
  wailsjs/go/service/     — Auto-generated TS bindings. Do not edit.
```

**Data flow:** Frontend → generated TS bindings → Go service methods → repository → PostgreSQL. No REST API.

**Core entities:** Zone, Client, Product, ClientProduct. Balance logic: `ClientProduct.balance += credit_units` on delivery confirmation. This is the only mutation path.

---

## Critical Patterns

### Wails binding + Init() lifecycle

Wails needs non-nil service pointers at `Bind` time (before `startup()`), but DB/repos aren't available yet. The pattern:

1. `NewApp()` allocates services as empty structs: `&service.ZoneService{}`
2. `startup()` connects DB, runs migrations/seeds, creates repos, then calls `Init(ctx, ...)` on each service

**If you add a new service:** allocate it in `NewApp()`, call `Init()` in `startup()`, add it to the `Bind` slice in `main.go`.

### DB connection

Env var `SODITAPP_DB_URL` (default: `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`).

Driver: **pgx/v5** via `database/sql` (`_ "github.com/jackc/pgx/v5/stdlib"`). Use `$1`-style params.

### Custom migration tracker

Not golang-migrate. The `RunMigrations()` function reads embedded `.sql` files, tracks applied versions in `schema_migrations`. Seeds use identical logic with `schema_seeds`. Both run automatically on app startup.

---

## Svelte 5 Rules

- **Runes only.** `$state()`, `$derived()`, `$effect()`, `$props()`, `$bindable()`. No `$:` or Svelte 4 stores.
- **Snippets over slots.** Use `{#snippet}` / `{@render}` — `<slot>` is deprecated.
- **Path aliases:** `$lib` → `src/lib/`, `wailsjs` → `frontend/wailsjs/`
- **shadcn-svelte:** style `mira`, icon library `hugeicons`. Components managed by CLI — don't hand-edit `src/lib/components/ui/`.

## Go Rules

- Package names: singular, lowercase (`model`, `repository`, `service`, `dto`).
- Repository layer owns all SQL. Services must not contain raw SQL.
- Wails-bound methods return `(result, error)` — frontend sees promise rejections.
- Keep `main.go` minimal — only bootstrap and embed directives.

## General

- **UI language: Spanish.** All user-facing text in Spanish; code in English.
- **Package manager:** `pnpm` for frontend.
- **No over-engineering.** Single-user desktop app. Keep it simple.
- **Entities are fixed:** Zone, Client, Product, ClientProduct. Don't add new entities without discussion.
- **DTOs are ephemeral:** `DeliverySheetDTO` / `DeliveryEntryDTO` are runtime-only for print preview.
