CREATE TABLE IF NOT EXISTS notifications(
    id serial not null unique ,
    title varchar(255) not null,
    partner_id int references partners(id) on delete cascade,
    text text,
    link text,
    status int check (notifications.status >= 1 and 2 >= notifications.status),
    noty_getters int check ( notifications.noty_getters >= 1 and 2 >= notifications.noty_getters ),
    reference text,
    noty_date timestamp,
    noty_time int
);

CREATE TABLE IF NOT EXISTS getters(
    id serial not null unique ,
    notification_id int references notifications(id) on delete cascade,
    user_id int references users(id) on delete cascade
);

