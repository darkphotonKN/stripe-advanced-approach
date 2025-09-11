-- Add stripe_customer_id column to payments table
ALTER TABLE payments 
ADD COLUMN IF NOT EXISTS stripe_customer_id VARCHAR(255);

-- Create index for stripe_customer_id for faster queries
CREATE INDEX IF NOT EXISTS idx_payments_stripe_customer_id ON payments(stripe_customer_id);

-- Also rename stripe_payment_intent_id to stripe_intent_id to match the code
ALTER TABLE payments 
RENAME COLUMN stripe_payment_intent_id TO stripe_intent_id;

-- Drop the old index and create new one with the new column name
DROP INDEX IF EXISTS idx_payments_stripe_payment_intent_id;
CREATE INDEX IF NOT EXISTS idx_payments_stripe_intent_id ON payments(stripe_intent_id);
