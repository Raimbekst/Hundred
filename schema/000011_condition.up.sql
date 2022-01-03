CREATE TABLE IF NOT EXISTS conditions(
    id serial not null unique ,
    caption varchar(255),
    text text
);