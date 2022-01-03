CREATE TYPE raffle_status_privileges AS ENUM ('запланирован','состоялся','не состоялся');

CREATE TABLE IF NOT EXISTS raffles(
   id serial not null unique ,
   raffle_date date not null,
   raffle_time int not null,
   check_category int not null,
   raffle_type int check ( raffles.raffle_type >= 1 and 3 >= raffles.raffle_type ),
   status raffle_status_privileges,
   reference text,
   check_id int references checks(id) on delete cascade
);

