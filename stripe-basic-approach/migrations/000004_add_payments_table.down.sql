-- Drop indexes
DROP INDEX IF EXISTS idx_payments_user_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_stripe_payment_intent_id;
DROP INDEX IF EXISTS idx_payments_stripe_session_id;

-- Drop payments table
DROP TABLE IF EXISTS payments;