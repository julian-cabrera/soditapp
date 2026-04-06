> [!SUMMARY] Delivery and Client management Application.

# Stack
Application: Wails (https://wails.io/)
- GO
- Svelte 5 with shadcn style components
Database: PostgreSQL
OS: Windows 10 or above
***
# Summary
## Objective
AS A Software Engineer with expertise on Java and React, I WANT to design a desktop application SO THAT I can develop it later.
## Context
A User has an Old software for his Delivery Business that can no longer be supported and does not run in Windows 10 or higher.
I've been asked to create a new software in replacement so that the User can update his computer and narrow down the features of the old software to his actual needs.
The application's purpose is to help the User handle Delivery and Clients information for it's delivery business, keep it stored and available at any moment.

### Clarified Requirements
- No money/payments — only **product units** are tracked (charged vs. on credit)
- 2 products currently, but model supports more
- Units are entered as raw counts — no unit conversion
- Any client can receive any product
- Delivery Sheets have **4 generic day columns**, organized by Zone
- Backups are manual (user-triggered export)

### Daily routine use case example
1. The user delivers Products for many Clients.
2. When he makes his deliveries, the user writes down on a paper Delivery Sheet how many units he was paid for or how many units were sold on credit.
3. When the delivery is complete, the User uploads the information manually to the Application that stores said information.
4. The Application keeps track of each Client's products and calculates their balance of units owed.
### Delivery management
Deliveries sheets are organized based on City Zones
Each Client belongs to a certain Zone
Daily sheets with Clients are organized by Zone
Each sheet has 4 generic day columns — the user handwrites the actual dates on paper
#### Delivery Sheet example

The sheet always has exactly **4 generic `Day` columns**. The user handwrites the actual dates on paper.
Only clients with a **negative balance** appear. If the client list overflows a single A4 page, additional pages are printed with the same layout and an incremented page number.

```
Zone: {zone_name}          Date: {created_at}          Page: {page_number}

Client Info          | Day         | Day         | Day         | Day
                     | Chg | Cred  | Chg | Cred  | Chg | Cred  | Chg | Cred
─────────────────────┼─────┼───────┼─────┼───────┼─────┼───────┼─────┼──────
{name} - {bal}       |     |       |     |       |     |       |     |
  {address}          |     |       |     |       |     |       |     |
```

- **Charged (Chg):** units paid on the spot — do not affect balance
- **On Credit (Cred):** units to be paid in the future — increase the client's debt

#### Entities

###### Zone
```
id          SERIAL PRIMARY KEY
name        VARCHAR NOT NULL
```

###### Client
```
id                SERIAL PRIMARY KEY
name              VARCHAR NOT NULL
address           VARCHAR
zone_id           FK -> Zone
create_date       DATE NOT NULL DEFAULT CURRENT_DATE
unregister_date   DATE
```

###### Product
```
id          SERIAL PRIMARY KEY
name        VARCHAR NOT NULL
active      BOOLEAN NOT NULL DEFAULT TRUE
```

###### ClientProduct (tracks running credit balance per client per product)
```
id          SERIAL PRIMARY KEY
client_id   FK -> Client
product_id  FK -> Product
balance     INTEGER NOT NULL DEFAULT 0   -- units currently owed
UNIQUE(client_id, product_id)
```

#### Balance Logic
- When a delivery entry is confirmed: `ClientProduct.balance += credit_units`
- When a delivery entry is edited/deleted: recalculate delta accordingly

#### DTOs (Runtime Only)
`DeliverySheet` and `DeliveryEntry` are **not persisted**. They are Go structs built at runtime for print preview and discarded after printing.

###### DeliverySheetDTO
```
zone_id       INTEGER       -- reference to Zone
zone_name     VARCHAR
created_at    TIMESTAMP     -- date the sheet is generated
page_number   INTEGER       -- 1-based; increments when client list overflows a page
entries       []DeliveryEntryDTO
```

###### DeliveryEntryDTO
```
client_id       INTEGER
client_name     VARCHAR
client_address  VARCHAR
client_balance  INTEGER     -- current balance (only negative values included)
product_id      INTEGER
product_name    VARCHAR
charged_units   INTEGER
credit_units    INTEGER
```

## Key MUST HAVE Features
- CRUD of Clients/Zones/Products
- Print Delivery Sheet into a A4 sheet.
- Import/Export CSV
- Rearrange Zones by removing or adding a group of Clients
- Delivery Sheet preview before printing

## MVP Breakdown

### MVP 1 — Core Data & Delivery Entry

#### Task 1 — Project Setup ✅
- Initialize Wails project (Go backend + Svelte 5 frontend)
- Configure PostgreSQL connection (connection string, env config)
- Set up folder structure: `backend/`, `frontend/`, `db/migrations/`
- Verify dev and build pipeline runs clean

#### Task 2 — Backend: Database & Entities ✅
- Write SQL migrations for: `Zone`, `Client`, `Product`, `ClientProduct`
- Implement repository layer with CRUD for each entity
- Implement `DeliverySheetDTO` / `DeliveryEntryDTO` builder (queries `ClientProduct` for negative balances, paginates by page)

#### Task 3 — Backend: Services & Wails Bindings ✅
- Zone service (create, list, update, delete)
- Client service (create, list by zone, update, unregister)
- Product service (create, list active, toggle active)
- ClientProduct service (ensure record exists, adjust balance on delivery confirm/edit/delete)
- Expose all services as Wails-bound methods callable from the frontend

#### Task 4 — Frontend: Zones & Clients ✅
- Zones list screen — create, rename, delete
- Clients list screen (filterable by zone) — create, edit, unregister
- Assign/reassign client to zone

#### Task 5 — Frontend: Products & Balances ✅
- Products list screen — create, rename, toggle active
- ClientProduct balance view per client (read-only, shows current balance per product)

#### Task 6 — Frontend: Delivery Entry ✅
- Select zone → view clients with negative balance
- Input charged/credit units per client per product
- On confirm: call ClientProduct service to update balances

### MVP 2 — Printing & Sheet Preview

> All tasks are frontend-only. The backend `GetDeliverySheets(zoneID)` already provides all needed data.
> Print strategy: `window.print()` + `@media print` CSS (Wails v2 has no dedicated print API).

#### Task 1 — Week Selector Component ✅
- New `WeekSelector.svelte` in `frontend/src/components/delivery/`
- Single date input → derives Monday–Sunday range automatically
- Add `weekStart`/`weekEnd` state to `delivery.svelte.ts`; default to current week
- Label in Spanish: "Semana del {start} al {end}" — cosmetic only, does not filter backend data
- **Files:** `WeekSelector.svelte` (new), `delivery.svelte.ts`, `DeliveryPage.svelte`

#### Task 2 — Delivery Sheet Preview Component ✅
- New `DeliverySheetPreview.svelte` in `frontend/src/components/delivery/`
- Props: `DeliverySheetDTO[]` + week label
- Per sheet: header (Zona, Semana, Página) + table matching the documented 4-day-column format
- 4 generic "Día" columns, each with Cob/Fiado sub-columns — cells left empty for handwriting
- Entries grouped by client: name + balance on row 1, address on row 2
- Plain HTML `<table>` (not shadcn Table) for precise print control; purely presentational
- **Files:** `DeliverySheetPreview.svelte` (new)

#### Task 3 — Print Stylesheet & Print Button ✅
- Add `@page { size: A4; margin: 10mm; }` and `@media print` rules in `app.css`
- Hide app shell (nav, controls) in print media; `page-break-after: always` between sheets
- "Imprimir" button calls `window.print()`; add `print:hidden` to `AppShell.svelte` nav
- **Files:** `app.css`, `DeliveryPage.svelte`, `DeliverySheetPreview.svelte`, `AppShell.svelte`

#### Task 4 — Integration: Wire Preview into Delivery Page ✅
- "Ver Hoja" button toggles from entry form → preview mode (week selector + sheet preview + print button)
- "Volver" returns to entry form preserving all edits and state
- Extract `groupedEntries` logic to shared utility reused by both entry form and preview
- Add `previewMode` state to `delivery.svelte.ts`
- **Files:** `DeliveryPage.svelte`, `delivery.svelte.ts`

**Task dependencies:** Tasks 1–3 are independent. Task 4 depends on all three.

### MVP 3 — Data Management ✅

> Backend changes for export (file I/O via Wails runtime dialogs), import (CSV parsing + validation), and bulk reassignment (transaction wrapping existing repo). Search/filtering is frontend-only.

#### Task 1 — Backend: Export Service ✅
- New `backend/service/export_service.go` with `ExportService` struct bound in `app.go`
- Three methods, each opens a Wails `runtime.SaveFileDialog`, writes the CSV, returns error on cancel or failure:
  - `ExportZones()` → columns: `id, nombre`
  - `ExportClients()` → columns: `id, nombre, direccion, zona`
  - `ExportBalances()` → columns: `cliente, producto, saldo` (query joins `client_product` → `client` → `product`)
- Reuse existing repo `GetAll` methods; add a join query for balances to `ClientProductRepository`
- **Files:** `backend/service/export_service.go` (new), `backend/repository/client_product_repository.go`, `app.go`

#### Task 2 — Backend: Import Service ✅
- New `backend/service/import_service.go` with `ImportService` struct bound in `app.go`
- Single method `ImportClients()`:
  - Opens `runtime.OpenFileDialog` (filter `*.csv`)
  - Parses rows: `nombre` (required), `direccion` (optional), `zona` (must match an existing zone name)
  - Inserts valid rows, skips duplicates by name+zone, collects per-row errors
  - Returns `ImportResult { Inserted, Skipped, Errors []string }`
- **Files:** `backend/service/import_service.go` (new), `app.go`

#### Task 3 — Backend: Bulk Zone Reassignment ✅
- Add `BulkUpdateZone(ctx, clientIDs []int, targetZoneID int) error` to `ClientRepository` — single `UPDATE … WHERE id = ANY($1)` inside a transaction
- Add `MoveClientsToZone(clientIDs []int, targetZoneID int) error` to `ClientService` delegating to the new repo method
- **Files:** `backend/repository/client_repository.go`, `backend/service/client_service.go`

#### Task 4 — Frontend: Export Buttons ✅
- "Exportar" button on `ZonesPage.svelte` → calls `ExportZones()`
- "Exportar" button on `ClientsPage.svelte` → calls `ExportClients()` and `ExportBalances()` (two buttons or a dropdown)
- Show a sonner toast on success or error
- **Files:** `ZonesPage.svelte`, `ClientsPage.svelte`

#### Task 5 — Frontend: Import Dialog ✅
- New `ImportClientsDialog.svelte` — triggered by an "Importar Clientes" button in `ClientsPage.svelte`
- Calls `ImportClients()`, displays result summary (insertados / omitidos / errores) in the dialog
- On success, refreshes the client list
- **Files:** `frontend/src/components/clients/ImportClientsDialog.svelte` (new), `ClientsPage.svelte`

#### Task 6 — Frontend: Bulk Zone Reassignment UI ✅
- New `ZoneReassignDialog.svelte` triggered from `ZonesPage.svelte`
- Step 1: select source zone and target zone from dropdowns
- Step 2: checklist of active clients in the source zone (all pre-checked)
- Confirm calls `MoveClientsToZone`, refreshes zones and clients stores
- **Files:** `frontend/src/components/zones/ZoneReassignDialog.svelte` (new), `ZonesPage.svelte`

#### Task 7 — Frontend: Search & Text Filtering ✅
- Text search input on `ClientsPage.svelte` — client-side filter by name or address (case-insensitive, applied after zone filter)
- Text search input on `ZonesPage.svelte` — filter zone list by name
- No backend changes needed (data sets are small enough for client-side filtering)
- **Files:** `ClientsPage.svelte`, `ZonesPage.svelte`

**Task dependencies:** Tasks 1–3 are independent backend tasks. Task 4 depends on Task 1. Task 5 depends on Task 2. Task 6 depends on Task 3. Task 7 is independent.

### MVP 4 — Reporting

> Requires a new `delivery_transaction` table persisted on every `ConfirmDelivery` call. Tasks 1–3 are backend; Tasks 4–5 are frontend. Task 1 → Task 2 → Task 3 → Tasks 4 & 5.

#### Task 1 — Backend: Delivery Transaction Table & Repository
- Migration: `delivery_transaction` (id, client_id FK → Client, product_id FK → Product, charged_units INTEGER, credit_units INTEGER, delivery_date DATE DEFAULT CURRENT_DATE)
- Model: `DeliveryTransaction` struct in `backend/model/`
- Repository: `backend/repository/delivery_transaction_repository.go`
  - `InsertTx(ctx, tx, clientID, productID, chargedUnits, creditUnits)` — called inside ConfirmDelivery DB transaction
  - `GetByClient(ctx, clientID, startDate, endDate string)` → `[]DeliveryTransactionRow`
  - `GetByZone(ctx, zoneID, startDate, endDate string)` → `[]DeliveryTransactionRow` (JOINs client to filter by zone)
- **Files:** `db/migrations/005_delivery_transaction.sql` (new), `backend/model/delivery_transaction.go` (new), `backend/repository/delivery_transaction_repository.go` (new)

#### Task 2 — Backend: Record Transactions on Delivery Confirm
- Inject `DeliveryTransactionRepository` into `ClientProductService`
- Extend `ConfirmDelivery` to call `InsertTx` within the existing DB transaction for each non-zero adjustment
- **Files:** `backend/service/client_product_service.go`, `app.go`

#### Task 3 — Backend: Report Service & DTOs
- New `backend/dto/report.go`:
  - `DeliveryHistoryEntryDTO` (clientID, clientName, productID, productName, chargedUnits, creditUnits, deliveryDate)
  - `ZoneReportEntryDTO` (same fields)
- New `backend/service/report_service.go` with `ReportService` bound in `app.go`:
  - `GetClientHistory(clientID int, startDate, endDate string) ([]DeliveryHistoryEntryDTO, error)`
  - `GetZoneReport(zoneID int, startDate, endDate string) ([]ZoneReportEntryDTO, error)`
- **Files:** `backend/dto/report.go` (new), `backend/service/report_service.go` (new), `app.go`

#### Task 4 — Frontend: Client Delivery History Dialog
- New `DeliveryHistoryDialog.svelte` — props: `clientId`, `clientName`, `open` ($bindable)
- On open: start/end date inputs (default last 30 days) + "Buscar" button → calls `GetClientHistory`
- Table columns: Fecha | Producto | Cobrado | Fiado
- Triggered by "Ver Historial" button added to `ClientsPage.svelte`
- **Files:** `frontend/src/components/clients/DeliveryHistoryDialog.svelte` (new), `frontend/src/components/clients/ClientsPage.svelte`

#### Task 5 — Frontend: Zone Reports Page
- New `ZoneReportsPage.svelte` — zone selector + start/end date inputs + "Buscar" button
- Table: Fecha | Cliente | Producto | Cobrado | Fiado, sorted by date then client
- Add "Reportes" nav entry in `AppShell.svelte`; wire page in `App.svelte`
- **Files:** `frontend/src/components/zones/ZoneReportsPage.svelte` (new), `frontend/src/components/AppShell.svelte`, `frontend/src/App.svelte`
