CREATE TABLE IF NOT EXISTS courier (
    id serial PRIMARY KEY,
    courier_type varchar(50) NOT NULL,
    regions integer[] NOT NULL,
    working_hours text[] NOT NULL
);

CREATE TABLE IF NOT EXISTS "order" (
    id serial PRIMARY KEY,
    courier_id integer REFERENCES courier,
    weight integer NOT NULL,
    region integer NOT NULL,
    delivery_time varchar(255),
    complete_time timestamp,
    price integer NOT NULL
);
