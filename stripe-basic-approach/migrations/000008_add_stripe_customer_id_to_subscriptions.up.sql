-- Add stripe_customer_id to subscriptions table
ALTER TABLE subscriptions ADD COLUMN stripe_customer_id VARCHAR(255);

-- Create index for common queries by customer ID
CREATE INDEX idx_subscriptions_stripe_customer_id ON subscriptions(stripe_customer_id);