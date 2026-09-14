CREATE TABLE IF NOT EXISTS delivery_transaction (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id       TEXT NOT NULL REFERENCES client(reference_uuid),
    product_id      TEXT NOT NULL REFERENCES product(reference_uuid),
    charged_units   INTEGER NOT NULL DEFAULT 0,
    credit_units    INTEGER NOT NULL DEFAULT 0,
    delivery_date   TEXT NOT NULL DEFAULT (date('now'))
);
