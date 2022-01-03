CREATE TYPE status_noty_privileges AS ENUM ('отправлено','запланировано');
CREATE TYPE get_notification_privileges AS ENUM ('все','выборочно');
CREATE TABLE IF NOT EXISTS notifications(
    id serial not null unique ,
    title varchar(255) not null,
    partner_id int references partners(id) on delete cascade,
    text text,
    link text,
    status status_noty_privileges,
    noty_getters get_notification_privileges,
    reference text,
    noty_date date,
    noty_time float
);

CREATE TABLE IF NOT EXISTS getters(
    id serial not null unique ,
    notification_id int references notifications(id) on delete cascade,
    user_id int references users(id) on delete cascade
);

