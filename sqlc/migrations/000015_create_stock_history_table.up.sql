CREATE TYPE stock_change_type AS ENUM ('sale', 'restock', 'correction', 'return', 'damage');

CREATE TABLE stock_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    change_amount INTEGER NOT NULL,
    previous_stock INTEGER NOT NULL,
    current_stock INTEGER NOT NULL,
    change_type stock_change_type NOT NULL,
    reference_id UUID,
    note TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_history_product_id ON stock_history(product_id);
CREATE INDEX idx_stock_history_created_at ON stock_history(created_at);
