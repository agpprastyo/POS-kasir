ALTER TABLE orders ADD COLUMN tax_amount BIGINT NOT NULL DEFAULT 0 CHECK (tax_amount >= 0);
ALTER TABLE orders ADD COLUMN service_charge_amount BIGINT NOT NULL DEFAULT 0 CHECK (service_charge_amount >= 0);
