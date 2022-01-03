CREATE TABLE IF NOT EXISTS partners(
  id serial not null unique ,
  position int,
  partner_name varchar(255) not null,
  logo text ,
  link_website text,
  banner text ,
  status int check ( partners.status >= 1 and 2 >= partners.status),
  start_partnership varchar(255),
  end_partnership varchar(255),
  partner_package varchar(255),
  reference text
);