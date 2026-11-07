CREATE TABLE IF NOT EXISTS client (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    name              TEXT NOT NULL UNIQUE,
    address           TEXT,
    zone_id           TEXT NOT NULL REFERENCES zone(reference_uuid),
    create_date       TEXT NOT NULL DEFAULT (date('now')),
    unregister_date   TEXT,
    priority_order    INTEGER NOT NULL DEFAULT 0,
    reference_uuid    TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_client_reference_uuid ON client(reference_uuid);
