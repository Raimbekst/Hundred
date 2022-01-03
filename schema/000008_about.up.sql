CREATE TABLE IF NOT EXISTS about(
    id serial not null unique ,
    facebook_link text,
    youtube_link text,
    instagram_link text,
    tiktok_link text,
    whatsapp_link text,
    telegram_link text,
    phone_number varchar(255),
    phone_number_2 varchar(255)
);

CREATE TABLE  IF NOT EXISTS descriptions(
    id serial not null unique ,
    caption varchar(255),
    text text
);
