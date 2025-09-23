-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_user_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_stripe_subscription_id;

-- Drop subscriptions table
DROP TABLE IF EXISTS subscriptions;