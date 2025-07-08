-- ENUMS
CREATE TYPE user_role AS ENUM ('admin', 'cashier', 'manager');
CREATE TYPE order_type AS ENUM ('dine_in', 'takeaway');
CREATE TYPE order_status AS ENUM ('open', 'in_progress', 'served', 'paid', 'cancelled');
CREATE TYPE discount_type AS ENUM ('percentage', 'fixed_amount');
CREATE TYPE promotion_scope AS ENUM ('ORDER', 'ITEM');
CREATE TYPE promotion_target_type AS ENUM ('PRODUCT', 'CATEGORY');
CREATE TYPE promotion_rule_type AS ENUM (
  'MINIMUM_ORDER_AMOUNT',
  'REQUIRED_PRODUCT',
  'REQUIRED_CATEGORY',
  'ALLOWED_PAYMENT_METHOD',
  'ALLOWED_ORDER_TYPE'
);
CREATE TYPE log_action_type AS ENUM (
  'CREATE', 'UPDATE', 'DELETE', 'CANCEL', 'APPLY_PROMOTION', 'PROCESS_PAYMENT'
);
CREATE TYPE log_entity_type AS ENUM (
  'PRODUCT', 'CATEGORY', 'PROMOTION', 'ORDER', 'USER'
);

-- TABLES
CREATE TABLE users (
  id UUID PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  avatar TEXT,
  role user_role NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
  id UUID PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  image_url TEXT,
  price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
  stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_options (
  id UUID PRIMARY KEY,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  additional_price NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (additional_price >= 0),
  image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_methods (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cancellation_reasons (
  id SERIAL PRIMARY KEY,
  reason VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE promotions (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  scope promotion_scope NOT NULL DEFAULT 'ORDER',
  discount_type discount_type NOT NULL,
  discount_value NUMERIC(12,2) NOT NULL CHECK (discount_value >= 0),
  max_discount_amount NUMERIC(12,2) CHECK (max_discount_amount >= 0),
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE promotion_rules (
  id UUID PRIMARY KEY,
  promotion_id UUID NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
  rule_type promotion_rule_type NOT NULL,
  rule_value VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE promotion_targets (
  id UUID PRIMARY KEY,
  promotion_id UUID NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
  target_type promotion_target_type NOT NULL,
  target_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orders (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  type order_type NOT NULL DEFAULT 'dine_in',
  status order_status NOT NULL DEFAULT 'open',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  gross_total NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (gross_total >= 0),
  discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_total NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (net_total >= 0),
  applied_promotion_id UUID REFERENCES promotions(id) ON DELETE SET NULL,
  payment_method_id INTEGER REFERENCES payment_methods(id),
  payment_gateway_reference VARCHAR(255),
  cash_received NUMERIC(12,2) CHECK (cash_received >= 0),
  change_due NUMERIC(12,2) CHECK (change_due >= 0),
  cancellation_reason_id INTEGER REFERENCES cancellation_reasons(id) ON DELETE SET NULL,
  cancellation_notes TEXT
);

CREATE TABLE order_items (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
  price_at_sale NUMERIC(12,2) NOT NULL CHECK (price_at_sale >= 0),
  subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
  discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_subtotal NUMERIC(12,2) NOT NULL CHECK (net_subtotal >= 0)
);

CREATE TABLE order_item_options (
  id UUID PRIMARY KEY,
  order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
  product_option_id UUID NOT NULL REFERENCES product_options(id) ON DELETE RESTRICT,
  price_at_sale NUMERIC(12,2) NOT NULL CHECK (price_at_sale >= 0)
);

CREATE TABLE activity_logs (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action_type log_action_type NOT NULL,
  entity_type log_entity_type NOT NULL,
  entity_id VARCHAR(255) NOT NULL,
  details JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- INDEXES
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_promotions_dates ON promotions(start_date, end_date);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_payment_gateway_reference ON orders(payment_gateway_reference);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_entity ON activity_logs(entity_type, entity_id);


-- Automatically update updated_at on row update
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add triggers for tables with updated_at
DO $$
DECLARE
  tbl TEXT;
BEGIN
  FOR tbl IN
    SELECT table_name
    FROM information_schema.columns
    WHERE column_name = 'updated_at'
      AND table_schema = 'public'
  LOOP
    EXECUTE format('
      CREATE TRIGGER %I_set_updated_at
      BEFORE UPDATE ON %I
      FOR EACH ROW
      EXECUTE FUNCTION set_updated_at();
    ', tbl, tbl);
  END LOOP;
END;
$$;

-- Default UUID generation for id columns
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE products ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE product_options ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE promotions ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE promotion_rules ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE promotion_targets ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE orders ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE order_items ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE order_item_options ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE activity_logs ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Enable pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;