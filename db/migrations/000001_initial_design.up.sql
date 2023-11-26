CREATE TABLE "batches" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "activation_status" varchar NOT NULL,
  "customer_id" bigserial NOT NULL,
  "no_of_devices" integer NOT NULL,
  "mrc_id" bigserial NOT NULL,
  "delivery_date" timestamptz,
  "warranty_end" timestamptz
);

CREATE TABLE "bundles" (
  "id" bigserial PRIMARY KEY,
  "mrc" integer NOT NULL,
  "description" text NOT NULL
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "batch_id" bigserial NOT NULL,
  "bundle_id" bigserial NOT NULL,
  "nrc" boolean
);

CREATE TABLE "payment_logs" (
  "id" bigserial PRIMARY KEY,
  "payment" float NOT NULL,
  "due_date" timestamptz NOT NULL,
  "confirmation_date" timestamptz,
  "order_id" bigserial NOT NULL,
  "confirmed" boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE "charges" (
  "id" bigserial PRIMARY KEY,
  "paid" float NOT NULL,
  "due" float NOT NULL,
  "customer_id" bigserial NOT NULL
);

CREATE TABLE "customers" (
  "id" bigserial PRIMARY KEY,
  "customer" varchar UNIQUE NOT NULL
);

ALTER TABLE "orders" ADD FOREIGN KEY ("batch_id") REFERENCES "batches" ("id");

ALTER TABLE "payment_logs" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("bundle_id") REFERENCES "bundles" ("id");

ALTER TABLE "charges" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

CREATE TABLE "bundles_customers" (
  "bundles_id" bigserial,
  "customers_id" bigserial,
  PRIMARY KEY ("bundles_id", "customers_id")
);

ALTER TABLE "bundles_customers" ADD FOREIGN KEY ("bundles_id") REFERENCES "bundles" ("id");

ALTER TABLE "bundles_customers" ADD FOREIGN KEY ("customers_id") REFERENCES "customers" ("id");


ALTER TABLE "batches" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");
