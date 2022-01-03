CREATE TABLE IF NOT EXISTS checks(
  id serial not null unique ,
  user_id int references users(id) on delete cascade not null,
  partner_id int references partners(id) on delete cascade not null,
  check_amount int8 not null,
  check_date DATE not null,
  is_winner bool not null default false,
  registered_at timestamp with time zone default current_timestamp
);

CREATE TABLE IF NOT EXISTS check_images(
  id serial not null unique,
  check_id int references checks(id) on delete cascade not null ,
  check_image text not null
);
