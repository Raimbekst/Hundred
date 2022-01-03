CREATE TABLE IF NOT EXISTS banners(
  id serial not null unique,
  name varchar(255),
  status  int check ( banners.status >= 1 and 2 >= banners.status),
  image text,
  iframe text
);
