CREATE TYPE shift_status AS ENUM ('open', 'closed');

CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    start_cash BIGINT NOT NULL DEFAULT 0,
    expected_cash_end BIGINT,
    actual_cash_end BIGINT,
    status shift_status NOT NULL DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Ensure a user can have only one open shift at a time
CREATE UNIQUE INDEX idx_shifts_user_open ON shifts (user_id) WHERE status = 'open';

CREATE TYPE cash_transaction_type AS ENUM ('cash_in', 'cash_out');

CREATE TABLE cash_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shift_id UUID NOT NULL REFERENCES shifts(id),
    user_id UUID NOT NULL REFERENCES users(id),
    amount BIGINT NOT NULL,
    type cash_transaction_type NOT NULL,
    category VARCHAR(50) NOT NULL, -- e.g., "expense", "drop", "float"
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
