

CREATE TABLE IF NOT EXISTS capsules (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    user_id int
);

CREATE TABLE IF NOT EXISTS capsules_items (
    id SERIAL PRIMARY KEY,
    capsule_id INT NOT NULL,
    item_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
                                              id SERIAL PRIMARY KEY NOT NULL,
                                              user_id INT NOT NULL,
                                              name VARCHAR NOT NULL,
    url VARCHAR,
    category INT NOT NULL,
    size_number int,
    size_text varchar,
    description varchar,
    color integer,
    file_name varchar
);

CREATE TABLE IF NOT EXISTS looks (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    look_name varchar NOT NULL,
    stylist_id int
);


CREATE TABLE IF NOT EXISTS looks_items (
                                     look_id INT NOT NULL,
                                     item_id INT NOT NULL
);