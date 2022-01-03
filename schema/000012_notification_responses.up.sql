CREATE TABLE IF NOT EXISTS notification_users(
    id serial not null unique ,
    user_id int references users(id) on delete cascade,
    registration_token text not null
);