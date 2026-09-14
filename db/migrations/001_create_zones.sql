CREATE TABLE IF NOT EXISTS zone (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    reference_uuid  TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_zone_reference_uuid ON zone(reference_uuid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_zone_name ON zone(name);
