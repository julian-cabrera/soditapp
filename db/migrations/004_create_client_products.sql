CREATE TABLE IF NOT EXISTS client_product (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id   TEXT NOT NULL REFERENCES client(reference_uuid),
    product_id  TEXT NOT NULL REFERENCES product(reference_uuid),
    balance     INTEGER NOT NULL DEFAULT 0,
    UNIQUE(client_id, product_id)
);
