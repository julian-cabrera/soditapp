# SoditApp

> Delivery and Client Management Desktop Application.

## About

SoditApp helps a delivery business owner manage clients, zones, products, and unit-based credit balances. 
It was created as a specific need to replace legacy software that no longer runs on modern Windows versions.

**Key features:**
- CRUD for Clients, Zones, and Products
- Track per-client product balances (charged vs. on credit)
- Print delivery sheets (A4 format, organized by zone)
- CSV import/export
- Bulk zone reassignment

## Tech Stack

- **Desktop framework:** [Wails v2](https://wails.io/) (Go + embedded web frontend)
- **Backend:** Go 1.23
- **Frontend:** Svelte 5, TypeScript, Tailwind CSS v4, shadcn-svelte
- **Database:** PostgreSQL
- **OS:** Windows 10+

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- PostgreSQL instance

## Development

```bash
# Install frontend dependencies
cd frontend && pnpm install && cd ..

# Run in live development mode (hot reload for frontend, auto-rebuild for Go)
wails dev
```

A Vite dev server provides hot reload for frontend changes. A secondary dev server at `http://localhost:34115` allows calling Go methods from the browser devtools.

## Building

```bash
# Build a production executable
wails build
```

The output binary is placed in `build/bin/`.

## Project Structure

```
soditapp/
├── main.go              # Wails entry point
├── app.go               # App struct and lifecycle hooks
├── frontend/            # Svelte 5 + TypeScript frontend
│   └── src/
│       ├── App.svelte
│       ├── lib/components/ui/   # shadcn-svelte components
│       └── lib/utils.ts         # Tailwind class merge utility
├── docs/soditapp.md     # Requirements and design document
├── AGENTS.md            # AI agent guidelines (Go & Svelte 5 practices)
└── CLAUDE.md            # Claude Code configuration
```

## Documentation

- [Requirements & Design](docs/soditapp.md) — entities, balance logic, delivery sheet format, MVP breakdown
- [AGENTS.md](AGENTS.md) — Go and Svelte 5 best practices for AI-assisted development
