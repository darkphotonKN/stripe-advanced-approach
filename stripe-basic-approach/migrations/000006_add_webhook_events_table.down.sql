-- Drop indexes
DROP INDEX IF EXISTS idx_webhook_events_stripe_event_id;
DROP INDEX IF EXISTS idx_webhook_events_event_type;
DROP INDEX IF EXISTS idx_webhook_events_processed;
DROP INDEX IF EXISTS idx_webhook_events_created_at;

-- Drop webhook_events table
DROP TABLE IF EXISTS webhook_events;