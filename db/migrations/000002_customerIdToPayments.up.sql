ALTER TABLE payment_logs
ADD COLUMN customer_id bigserial NOT NULL;

ALTER TABLE "payment_logs" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");
