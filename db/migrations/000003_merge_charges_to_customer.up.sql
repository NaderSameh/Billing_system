ALTER TABLE customers
ADD COLUMN paid float NOT NULL;

ALTER TABLE customers
ADD COLUMN due float NOT NULL;

DROP TABLE IF EXISTS "charges" CASCADE;
