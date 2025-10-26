-- Drop the stripe_customer_id index
DROP INDEX IF EXISTS idx_payments_stripe_customer_id;

-- Remove stripe_customer_id column
ALTER TABLE payments
DROP COLUMN IF EXISTS stripe_customer_id;