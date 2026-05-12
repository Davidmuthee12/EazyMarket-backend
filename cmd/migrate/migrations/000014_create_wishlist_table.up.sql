CREATE TABLE IF NOT EXISTS wishlists(
  user_id    UUID REFERENCES users(uuid) ON DELETE CASCADE,
  product_id UUID REFERENCES products(id) ON DELETE CASCADE,
  created_at timestamp(0) with time zone NOT NULL DEFAULt now(),
  PRIMARY KEY (user_id, product_id) 
);