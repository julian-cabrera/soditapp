CREATE TABLE IF NOT EXISTS product (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    type            TEXT NOT NULL,
    active          INTEGER NOT NULL DEFAULT 1,
    reference_uuid  TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_reference_uuid ON product(reference_uuid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_name ON product(name);
