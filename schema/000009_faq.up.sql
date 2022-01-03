CREATE TABLE IF not exists faqs(
    id serial not null unique ,
    question varchar(255),
    answer text
);