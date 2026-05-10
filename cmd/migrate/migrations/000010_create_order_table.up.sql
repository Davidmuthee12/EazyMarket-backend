CREATE TABLE IF NOT EXISTS orders (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id          UUID NOT NULL REFERENCES users(uuid),
  status           VARCHAR(30) DEFAULT 'pending', -- pending | confirmed | shipped | delivered | cancelled | refunded
  total_amount     NUMERIC(12,2) NOT NULL,
  shipping_address JSONB,
  notes            TEXT,
  created_at       timestamp(0) with time zone NOT NULL DEFAULt now(),
  updated_at       timestamp(0) with time zone NOT NULL DEFAULt now()
);

CREATE TABLE IF NOT EXISTS order_items (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id       UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id     UUID NOT NULL REFERENCES products(id),
  vendor_id      UUID NOT NULL REFERENCES vendor_profiles(user_id),
  quantity       INT NOT NULL,
  unit_price     NUMERIC(12,2) NOT NULL,
  subtotal       NUMERIC(12,2) NOT NULL,
  commission_amt NUMERIC(12,2) NOT NULL,   -- platform cut
  vendor_payout  NUMERIC(12,2) NOT NULL    -- vendor receives
);
