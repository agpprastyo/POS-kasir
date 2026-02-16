-- Add RESTORE to log_action_type
ALTER TYPE log_action_type ADD VALUE 'RESTORE';

-- Add new entities to log_entity_type
ALTER TYPE log_entity_type ADD VALUE 'SETTINGS';
ALTER TYPE log_entity_type ADD VALUE 'SHIFT';
ALTER TYPE log_entity_type ADD VALUE 'PAYMENT_METHOD';
ALTER TYPE log_entity_type ADD VALUE 'CANCELLATION_REASON';
