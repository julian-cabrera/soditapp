CREATE TABLE IF NOT EXISTS client (
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(100) NOT NULL UNIQUE,
    address           VARCHAR(100),
    zone_id           INTEGER REFERENCES zone(id),
    create_date       DATE NOT NULL DEFAULT CURRENT_DATE,
    unregister_date   DATE,
    priority_order    INTEGER NOT NULL DEFAULT 0
);
