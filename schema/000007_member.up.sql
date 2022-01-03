CREATE TABLE IF NOT EXISTS members(
  id serial not null unique ,
  check_id int references checks(id) on delete cascade not null ,
  raffle_id int references raffles(id) on delete cascade not null,
  is_winner bool not null default false
);