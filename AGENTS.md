# AGENTS.md — SoditApp

## Project Overview

SoditApp is a desktop delivery and client management application built with **Wails v2** (Go backend + Svelte 5 frontend). It targets Windows 10+ and uses PostgreSQL for persistence.

The app helps a delivery business owner track clients, zones, products, and unit-based credit balances. There are no monetary transactions — only product unit counts (charged vs. on credit).

See `docs/soditapp.md` for full requirements, entity definitions, and MVP breakdown.

---

## Architecture

```
soditapp/
├── main.go                  # Wails entry point
├── app.go                   # App struct, lifecycle hooks
├── backend/                 # Go: services, repositories, models (to be created)
│   ├── models/              # Structs for Zone, Client, Product, ClientProduct
│   ├── repository/          # Database access layer (SQL queries)
│   ├── service/             # Business logic, Wails-bound methods
│   └── dto/                 # DeliverySheetDTO, DeliveryEntryDTO (runtime only)
├── db/migrations/           # SQL migration files (to be created)
├── frontend/
│   └── src/
│       ├── App.svelte       # Root component
│       ├── components/      # App-level components
│       ├── lib/
│       │   ├── components/ui/  # shadcn-svelte primitives
│       │   ├── hooks/       # Svelte 5 reusable hooks
│       │   └── utils.ts     # Utility functions (cn, etc.)
│       ├── app.css          # Global styles (Tailwind)
│       └── main.ts          # Frontend entry
├── wails.json               # Wails project config
└── docs/soditapp.md         # Requirements & design document
```

**Data flow:** Frontend (Svelte 5) calls Wails-bound Go methods directly via generated TypeScript bindings. Go services handle business logic and talk to PostgreSQL through the repository layer. `DeliverySheetDTO` / `DeliveryEntryDTO` are built at runtime for print preview and never persisted.

---

## Go Best Practices

### Project Structure
- Keep the `main.go` minimal — only Wails bootstrap.
- Place all domain logic under `backend/`. Organize by layer: `models/`, `repository/`, `service/`, `dto/`.
- Each entity (Zone, Client, Product, ClientProduct) gets its own file in each layer.

### Naming & Style
- Follow standard Go conventions: `MixedCaps` for exported names, `mixedCaps` for unexported.
- Use receiver names that are short and consistent (e.g., `(s *ZoneService)`, `(r *ClientRepo)`).
- Name interfaces by behavior: `ZoneReader`, `ClientWriter`, not `IZoneService`.
- Keep package names short, lowercase, singular: `model`, `repository`, `service`.

### Error Handling
- Always return `error` as the last return value from functions that can fail.
- Wrap errors with context using `fmt.Errorf("operation failed: %w", err)`.
- Never ignore errors silently. If a function returns an error, handle it or propagate it.
- Wails-bound methods should return `(result, error)` — the frontend receives these as promise rejections.

### Database
- Use `database/sql` with `lib/pq` or `pgx` for PostgreSQL.
- Keep SQL in the repository layer only. Services must not contain raw SQL.
- Use parameterized queries (`$1`, `$2`) — never string concatenation for query values.
- Migrations go in `db/migrations/` as numbered `.sql` files (e.g., `001_create_zones.sql`).
- Use transactions when multiple writes must succeed or fail together (e.g., delivery confirmation updating multiple `ClientProduct` balances).

### Wails Bindings
- Only expose service-level methods via Wails bindings — never repositories directly.
- Bound methods must be on the `App` struct or on structs registered in `Bind: []interface{}`.
- Keep bound method signatures simple: accept and return primitive types, slices, or plain structs. Avoid `interface{}` in bound method signatures.
- Use the Wails `runtime` package for events (e.g., notifying the frontend of state changes).

### Concurrency & Context
- Thread the `context.Context` from Wails startup through to database calls.
- Use `context.WithTimeout` for database operations to prevent hanging queries.
- Do not use goroutines for request handling — Wails already manages this.

### Testing
- Write table-driven tests for service logic.
- Use interfaces for repository dependencies to allow test doubles.
- Test files live next to the code they test (`zone_service_test.go` alongside `zone_service.go`).

---

## Svelte 5 Best Practices

### Runes & Reactivity
- Use Svelte 5 runes for all reactive state. Prefer `$state()` for local component state and `$derived()` for computed values.
- Never use the legacy `$:` reactive statement or Svelte 4 stores. This is a Svelte 5 project.
- Use `$effect()` sparingly and only for side effects (DOM manipulation, logging, external API calls). Do not use it for derived state — use `$derived()` instead.
- Use `$bindable()` for props that support two-way binding.
- Use `$props()` to declare component props: `let { name, value = 'default' } = $props()`.

### Component Patterns
- Keep components small and focused. If a component exceeds ~150 lines, consider splitting it.
- Use `{#snippet}` blocks for reusable template fragments within a component. Prefer snippets over splitting into tiny sub-components when the fragment is only used locally.
- Pass snippets as props (`children`, named snippets) instead of using `<slot>` (deprecated in Svelte 5).
- Use `{@render children()}` to render snippet props.

### TypeScript
- All `.svelte` files must use `<script lang="ts">`.
- Define explicit types for component props, Wails binding responses, and shared data structures.
- Place shared types in `frontend/src/lib/types/`.
- Use the Wails-generated TypeScript bindings — do not manually duplicate Go struct types.

### Styling
- Use Tailwind CSS utility classes for styling. The project uses Tailwind v4 with the Vite plugin.
- Use shadcn-svelte components from `$lib/components/ui/` for standard UI elements (buttons, inputs, dialogs, tables, etc.).
- Follow shadcn-svelte conventions: compose from primitives, customize via Tailwind classes and `class` prop, extend with `tailwind-variants` for complex variant logic.
- Use `cn()` from `$lib/utils` to merge class names.
- Keep `app.css` for global base styles only. Component-specific styles go in the component via Tailwind.

### State Management
- For global/shared state, use Svelte 5 module-level `$state()` in `.svelte.ts` files (e.g., `frontend/src/lib/stores/clients.svelte.ts`).
- Do not use Svelte 4 stores (`writable`, `readable`, `derived`). Use runes-based `.svelte.ts` modules instead.
- Keep state as close to where it's used as possible. Lift state only when needed.

### Wails Integration
- Call Go-bound methods via the generated bindings in `frontend/wailsjs/`.
- Handle async Go calls with `async/await`. Show loading states for operations that may take time.
- Handle errors from Go calls gracefully — display user-friendly messages, not raw error strings.

### File Organization
- `src/components/` — app-level composed components (pages, layouts, forms).
- `src/lib/components/ui/` — shadcn-svelte primitives (managed by shadcn CLI).
- `src/lib/hooks/` — reusable logic as Svelte 5 hooks (`.svelte.ts` files using runes).
- `src/lib/types/` — shared TypeScript interfaces and types.
- `src/lib/utils.ts` — small utility functions.

---

## General Rules

- **Language:** The application UI is in **Spanish**. All user-facing labels, messages, and text must be in Spanish. Code (variables, functions, comments) remains in English.
- **Package manager:** Use `pnpm` for the frontend (configured in `wails.json`).
- **No over-engineering:** This is a small single-user desktop app. Avoid unnecessary abstractions, DI frameworks, or enterprise patterns. Keep it simple and direct.
- **Entities are defined:** Zone, Client, Product, and ClientProduct are the core entities. Do not add new entities without explicit discussion.
- **Balance logic:** `ClientProduct.balance += credit_units` on delivery confirmation. Recalculate deltas on edit/delete. This is the only mutation path for balances.
- **DTOs are ephemeral:** `DeliverySheetDTO` and `DeliveryEntryDTO` are runtime-only Go structs for print preview. They are never persisted to the database.
- **Print layout:** Delivery sheets have exactly 4 generic day columns, organized by zone, showing all active clients. Pagination at A4 boundaries.
