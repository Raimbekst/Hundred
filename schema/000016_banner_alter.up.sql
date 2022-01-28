CREATE TYPE language_privileges AS ENUM ('kz','ru');
ALTER TABLE banners ADD COLUMN language_type language_privileges;