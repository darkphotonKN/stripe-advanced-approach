-- Remove stripe_customer_id from subscriptions table
DROP INDEX IF EXISTS idx_subscriptions_stripe_customer_id;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS stripe_customer_id;