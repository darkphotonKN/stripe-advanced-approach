-- Revert column name change
ALTER TABLE payments 
RENAME COLUMN stripe_intent_id TO stripe_payment_intent_id;

-- Drop the new index and recreate the old one
DROP INDEX IF EXISTS idx_payments_stripe_intent_id;
CREATE INDEX IF NOT EXISTS idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);

-- Drop the stripe_customer_id index
DROP INDEX IF EXISTS idx_payments_stripe_customer_id;

-- Remove stripe_customer_id column
ALTER TABLE payments 
DROP COLUMN IF EXISTS stripe_customer_id;