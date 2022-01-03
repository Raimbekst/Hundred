CREATE TYPE user_privileges AS ENUM ('admin','user');
CREATE TYPE gender_type  AS ENUM ('female','male');


SET TIMEZONE = 'Asia/Almaty';

CREATE TABLE IF NOT EXISTS users(
  id serial not null unique ,
  user_type user_privileges,
  user_name varchar(255) not null,
  phone_number varchar(255) not null,
  email varchar(255) not null,
  gender gender_type,
  age int,
  city varchar(255),
  password varchar(255) not null,
  registered_at timestamp with time zone default current_timestamp,
  is_blocked bool not null default false
);

CREATE TABLE IF NOT EXISTS sessions(
  id serial not null unique,
  user_id int references users(id) on delete cascade not null,
  refresh_token varchar(500) not null default '',
  expires_at timestamp with time zone default current_timestamp
);
