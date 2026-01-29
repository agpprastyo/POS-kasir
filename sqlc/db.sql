-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

CREATE TABLE public.activity_logs (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  action_type USER-DEFINED NOT NULL,
  entity_type USER-DEFINED NOT NULL,
  entity_id character varying NOT NULL,
  details jsonb,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT activity_logs_pkey PRIMARY KEY (id),
  CONSTRAINT activity_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE TABLE public.cancellation_reasons (
  id integer NOT NULL DEFAULT nextval('cancellation_reasons_id_seq'::regclass),
  reason character varying NOT NULL UNIQUE,
  description text,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT cancellation_reasons_pkey PRIMARY KEY (id)
);
CREATE TABLE public.categories (
  id integer NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
  name character varying NOT NULL UNIQUE,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT categories_pkey PRIMARY KEY (id)
);
CREATE TABLE public.order_item_options (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  order_item_id uuid NOT NULL,
  product_option_id uuid NOT NULL,
  price_at_sale bigint NOT NULL CHECK (price_at_sale::numeric >= 0::numeric),
  CONSTRAINT order_item_options_pkey PRIMARY KEY (id),
  CONSTRAINT order_item_options_order_item_id_fkey FOREIGN KEY (order_item_id) REFERENCES public.order_items(id),
  CONSTRAINT order_item_options_product_option_id_fkey FOREIGN KEY (product_option_id) REFERENCES public.product_options(id)
);
CREATE TABLE public.order_items (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL,
  product_id uuid NOT NULL,
  quantity integer NOT NULL DEFAULT 1 CHECK (quantity > 0),
  price_at_sale bigint NOT NULL CHECK (price_at_sale::numeric >= 0::numeric),
  subtotal bigint NOT NULL CHECK (subtotal::numeric >= 0::numeric),
  discount_amount bigint NOT NULL DEFAULT 0 CHECK (discount_amount::numeric >= 0::numeric),
  net_subtotal bigint NOT NULL CHECK (net_subtotal::numeric >= 0::numeric),
  CONSTRAINT order_items_pkey PRIMARY KEY (id),
  CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id),
  CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id)
);
CREATE TABLE public.orders (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  type USER-DEFINED NOT NULL DEFAULT 'dine_in'::order_type,
  status USER-DEFINED NOT NULL DEFAULT 'open'::order_status,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  gross_total bigint NOT NULL DEFAULT 0 CHECK (gross_total::numeric >= 0::numeric),
  discount_amount bigint NOT NULL DEFAULT 0 CHECK (discount_amount::numeric >= 0::numeric),
  net_total bigint NOT NULL DEFAULT 0 CHECK (net_total::numeric >= 0::numeric),
  applied_promotion_id uuid,
  payment_method_id integer,
  payment_gateway_reference character varying,
  cash_received bigint CHECK (cash_received::numeric >= 0::numeric),
  change_due bigint CHECK (change_due::numeric >= 0::numeric),
  cancellation_reason_id integer,
  cancellation_notes text,
  CONSTRAINT orders_pkey PRIMARY KEY (id),
  CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id),
  CONSTRAINT orders_applied_promotion_id_fkey FOREIGN KEY (applied_promotion_id) REFERENCES public.promotions(id),
  CONSTRAINT orders_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id),
  CONSTRAINT orders_cancellation_reason_id_fkey FOREIGN KEY (cancellation_reason_id) REFERENCES public.cancellation_reasons(id)
);
CREATE TABLE public.payment_methods (
  id integer NOT NULL DEFAULT nextval('payment_methods_id_seq'::regclass),
  name character varying NOT NULL UNIQUE,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT payment_methods_pkey PRIMARY KEY (id)
);
CREATE TABLE public.product_options (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL,
  name character varying NOT NULL,
  additional_price bigint NOT NULL DEFAULT 0 CHECK (additional_price::numeric >= 0::numeric),
  image_url text,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT product_options_pkey PRIMARY KEY (id),
  CONSTRAINT product_options_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id)
);
CREATE TABLE public.products (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying NOT NULL,
  category_id integer,
  image_url text,
  price bigint NOT NULL CHECK (price::numeric >= 0::numeric),
  stock integer NOT NULL DEFAULT 0 CHECK (stock >= 0),
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT products_pkey PRIMARY KEY (id),
  CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id)
);
CREATE TABLE public.promotion_rules (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  promotion_id uuid NOT NULL,
  rule_type USER-DEFINED NOT NULL,
  rule_value character varying NOT NULL,
  description text,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT promotion_rules_pkey PRIMARY KEY (id),
  CONSTRAINT promotion_rules_promotion_id_fkey FOREIGN KEY (promotion_id) REFERENCES public.promotions(id)
);
CREATE TABLE public.promotion_targets (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  promotion_id uuid NOT NULL,
  target_type USER-DEFINED NOT NULL,
  target_id character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT promotion_targets_pkey PRIMARY KEY (id),
  CONSTRAINT promotion_targets_promotion_id_fkey FOREIGN KEY (promotion_id) REFERENCES public.promotions(id)
);
CREATE TABLE public.promotions (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying NOT NULL,
  description text,
  scope USER-DEFINED NOT NULL DEFAULT 'ORDER'::promotion_scope,
  discount_type USER-DEFINED NOT NULL,
  discount_value numeric NOT NULL CHECK (discount_value >= 0::numeric),
  max_discount_amount numeric CHECK (max_discount_amount >= 0::numeric),
  start_date timestamp with time zone NOT NULL,
  end_date timestamp with time zone NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at timestamp with time zone,
  CONSTRAINT promotions_pkey PRIMARY KEY (id)
);
CREATE TABLE public.schema_migrations (
  version bigint NOT NULL,
  dirty boolean NOT NULL,
  CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);
CREATE TABLE public.users (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  username character varying NOT NULL UNIQUE,
  email character varying NOT NULL UNIQUE,
  password_hash character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  avatar text,
  role USER-DEFINED NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  deleted_at timestamp with time zone,
  refresh_token text,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);