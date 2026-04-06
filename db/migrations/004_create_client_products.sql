CREATE TABLE IF NOT EXISTS client_product (
    id          SERIAL PRIMARY KEY,
    client_id   INTEGER NOT NULL REFERENCES client(id),
    product_id  INTEGER NOT NULL REFERENCES product(id),
    balance     INTEGER NOT NULL DEFAULT 0,
    UNIQUE(client_id, product_id)
);
