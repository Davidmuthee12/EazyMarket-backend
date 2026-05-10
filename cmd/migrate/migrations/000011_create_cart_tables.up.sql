-- CART (can also live purely in Redis)
CREATE TABLE IF NOT EXISTS carts (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID UNIQUE REFERENCES users(uuid) ON DELETE CASCADE,
  created_at timestamp(0) with time zone NOT NULL DEFAULt now(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULt now()
);

CREATE TABLE IF NOT EXISTS cart_items (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id    UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id),
  quantity   INT NOT NULL DEFAULT 1,
  UNIQUE(cart_id, product_id)
);