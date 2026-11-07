# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Live development (Go backend + Svelte frontend with hot reload)
wails dev

# Production build (outputs redistributable binary)
wails build

# Frontend only — from frontend/ directory
pnpm install          # install dependencies
pnpm run dev          # vite dev server only (no Go backend)
pnpm run build        # production frontend build
pnpm run check        # svelte-check + TypeScript type checking

# Add shadcn-svelte components
pnpm dlx shadcn-svelte@latest add <component-name>
```

## Architecture

**Wails v2 desktop app:** Go backend exposes methods that Svelte calls directly via auto-generated TypeScript bindings (no REST API). PostgreSQL for persistence.

- `main.go` — Wails bootstrap, embeds `frontend/dist` via `//go:embed`
- `app.go` — `App` struct with Wails lifecycle hooks; structs registered in `Bind` become callable from JS
- `frontend/` — Svelte 5 + TypeScript + Tailwind v4 + shadcn-svelte
  - `src/lib/components/ui/` — shadcn-svelte primitives (managed by shadcn CLI, avoid manual edits)
  - `src/lib/utils.ts` — `cn()` helper for Tailwind class merging
  - Path alias: `$lib` → `src/lib/`
  - Reference: https://www.shadcn-svelte.com/llms.txt

### Data Flow
Frontend → Wails-generated TS bindings → Go bound methods → Service layer → Repository layer → PostgreSQL

### Core Entities
Zone, Client, Product, ClientProduct (running credit balance per client per product). `DeliverySheetDTO` and `DeliveryEntryDTO` are runtime-only Go structs for print preview — never persisted.

## Key Conventions

- **UI language is Spanish.** All user-facing text in Spanish; code in English.
- **Svelte 5 runes only.** Use `$state()`, `$derived()`, `$effect()`, `$props()`, `$bindable()`. No Svelte 4 stores or `$:` syntax.
- **Snippets over slots.** Use `{#snippet}` and `{@render}` instead of `<slot>` (deprecated in Svelte 5).
- **Package manager:** `pnpm` (configured in `wails.json`).
- **Frontend styling:** Tailwind utility classes; compose with shadcn-svelte components; use `cn()` for class merging.
- See `AGENTS.md` for detailed Go and Svelte 5 best practices.
- See `docs/soditapp.md` for full requirements, entity schemas, balance logic, and MVP breakdown.
