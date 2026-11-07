# SoditApp

> Delivery and Client Management Desktop Application.

## About

SoditApp helps a delivery business owner manage clients, zones, products, and unit-based credit balances.
It was created as a specific need to replace legacy software that no longer runs on modern Windows versions.

**Key features:**
- CRUD for Zones, Clients, and Products
- Track per-client product balances (charged vs. on credit)
- Print delivery sheets (A4 format, organized by zone, paginated)
- Confirm deliveries: adjust balances and record charged/credit units per transaction
- Delivery history per client and zone reports by date range
- Bulk zone reassignment and client reordering within a zone
- Soft-delete (unregister) and restore clients
- SQLite backup export/import
- Dark mode

## Tech Stack

- **Desktop framework:** [Wails v2](https://wails.io/) (Go + embedded web frontend)
- **Backend:** Go 1.25
- **Frontend:** Svelte 5, TypeScript, Tailwind CSS v4, shadcn-svelte (hugeicons)
- **Database:** SQLite (embedded via [modernc.org/sqlite](https://gitlab.com/cznic/sqlite), no external DB server)
- **OS:** Windows 10+

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 20.19+](https://nodejs.org/) (or 22.12+)
- [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

No database server is required. On first run the app creates its SQLite database in the OS user config directory (e.g. `%AppData%\SoditApp\soditapp.db` on Windows), then applies embedded migrations and seeds. Override the connection string with the `SODITAPP_DB_URL` environment variable.

## Development

```bash
# Install frontend dependencies
cd frontend && pnpm install && cd ..

# Run in live development mode (hot reload for frontend, auto-rebuild for Go)
wails dev
```

A Vite dev server provides hot reload for frontend changes. A secondary dev server at `http://localhost:34115` allows calling Go methods from the browser devtools.

Useful frontend scripts (run from `frontend/`):

```bash
pnpm run dev      # vite only (no Go backend)
pnpm run build    # vite production build
pnpm run check    # svelte-check + TypeScript type checking
```

## Building

```bash
# Build a production executable
wails build
```

The output binary is placed in `build/bin/`.

## Project Structure

```
soditapp/
├── main.go                  # Wails entry point (embeds frontend/dist, migrations, seeds)
├── app.go                   # App struct, service wiring, startup/shutdown lifecycle
├── backend/
│   ├── database/            # Connect(), RunMigrations(), RunSeeds()
│   ├── model/               # Domain structs (Zone, Client, Product, ClientProduct, DeliveryTransaction)
│   ├── repository/          # SQL access layer (one repo per entity)
│   ├── service/             # Business logic, bound to Wails
│   └── dto/                 # Runtime-only DTOs (delivery sheets, reports)
├── db/
│   ├── migrations/          # Numbered .sql schema files (embedded at build time)
│   └── seeds/               # Numbered .sql seed files (embedded at build time)
├── frontend/                # Svelte 5 + TypeScript frontend
│   └── src/
│       ├── components/      # Pages and dialogs, grouped by entity (zones/, clients/, delivery/, ...)
│       ├── lib/components/ui/   # shadcn-svelte components (managed by CLI)
│       ├── lib/stores/      # Global state as .svelte.ts files (Svelte 5 runes)
│       ├── lib/utils.ts     # Tailwind class merge utility
│       └── wailsjs/         # Auto-generated TypeScript bindings (do not edit)
├── AGENTS.md                # AI agent guidelines (Go & Svelte 5 practices)
└── CLAUDE.md                # Claude Code configuration
```

## Documentation

- [AGENTS.md](AGENTS.md) — Go and Svelte 5 best practices for AI-assisted development
- [CLAUDE.md](CLAUDE.md) — Claude Code configuration
