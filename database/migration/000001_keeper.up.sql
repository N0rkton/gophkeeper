BEGIN;

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL
    );
CREATE TABLE IF NOT EXISTS keeper (
    id SERIAL PRIMARY KEY,
    data_id varchar(255) NOT NULL ,
    user_id int references users(id) NOT NULL,
    data_info varchar(255) NOT NULL,
    meta_info varchar(255),
    changed_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted bool default false,
    UNIQUE (user_id, data_id)
    );
COMMIT;

