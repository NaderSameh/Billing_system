ALTER TABLE IF EXISTS "payment_logs" DROP CONSTRAINT IF EXISTS "payments_logs_customer_id_fkey";

ALTER TABLE "payment_logs" DROP COLUMN "customer_id";