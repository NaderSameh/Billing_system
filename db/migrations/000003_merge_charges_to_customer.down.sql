ALTER TABLE customers
DROP COLUMN paid,
DROP COLUMN due;



CREATE TABLE "charges" (
  "id" bigserial PRIMARY KEY,
  "paid" float NOT NULL DEFAULT 0,
  "due" float NOT NULL DEFAULT 0,
  "customer_id" bigserial UNIQUE NOT NULL
);