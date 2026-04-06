CREATE TABLE IF NOT EXISTS client (
    id                SERIAL PRIMARY KEY,
    name              VARCHAR NOT NULL,
    address           VARCHAR,
    zone_id           INTEGER REFERENCES zone(id),
    create_date       DATE NOT NULL DEFAULT CURRENT_DATE,
    unregister_date   DATE
);
