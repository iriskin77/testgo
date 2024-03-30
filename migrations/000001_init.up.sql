CREATE TABLE locations
(
    id SERIAL UNIQUE NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    zip INT NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cars
(
    id serial unique not null,
    unique_number varchar(255) not null,
    car_name VARCHAR(255) NOT NULL,
    load_capacity INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    car_location_id int REFERENCES locations (id) ON DELETE SET NULL
);

CREATE TABLE cargos (
    id serial unique not null,
    cargo_name varchar(255) not null,
    weight INT NOT NULL,
    description varchar(1024) not null,
    pick_up_location_id int REFERENCES locations (id) ON DELETE SET NULL,
    delivery_location_id int REFERENCES locations (id) ON DELETE SET NULL
);

CREATE TABLE file (
    id serial unique not null,
    name varchar(255) not null,
    path varchar(2048) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)