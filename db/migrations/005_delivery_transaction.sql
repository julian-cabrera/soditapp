CREATE TABLE IF NOT EXISTS delivery_transaction (
    id              SERIAL PRIMARY KEY,
    client_id       INTEGER NOT NULL REFERENCES client(id),
    product_id      INTEGER NOT NULL REFERENCES product(id),
    charged_units   INTEGER NOT NULL DEFAULT 0,
    credit_units    INTEGER NOT NULL DEFAULT 0,
    delivery_date   DATE NOT NULL DEFAULT CURRENT_DATE
);
