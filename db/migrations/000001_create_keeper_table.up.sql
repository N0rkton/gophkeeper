BEGIN;
#DROP TYPE IF EXISTS order_state;
#create type order_state as enum ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED', 'NEW');
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS keeper (
    id int PRIMARY KEY UNIQUE,
    user_id int references users(id),
    data_info varchar(255) ,
    meta_info varchar(255),
    changed_at timestamp with time zone default now()
);
COMMIT;

