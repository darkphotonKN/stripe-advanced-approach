-- Add stripe_customer_id column to payments table
ALTER TABLE payments 
ADD COLUMN IF NOT EXISTS stripe_customer_id VARCHAR(255);

-- Create index for stripe_customer_id for faster queries
CREATE INDEX IF NOT EXISTS idx_payments_stripe_customer_id ON payments(stripe_customer_id);

