CREATE TABLE IF NOT EXISTS analytics_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  vendor_id UUID NOT NULL REFERENCES vendor_profiles(user_id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(uuid) ON DELETE SET NULL,
  session_id TEXT,
  event_type VARCHAR(50) NOT NULL,
  product_id UUID REFERENCES products(id) ON DELETE SET NULL,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (event_type IN (
    'storefront_view',
    'product_view',
    'add_to_cart',
    'wishlist_add',
    'checkout_started',
    'order_created'
  ))
);

CREATE INDEX IF NOT EXISTS analytics_events_vendor_created_idx
ON analytics_events(vendor_id, created_at DESC);

CREATE INDEX IF NOT EXISTS analytics_events_vendor_type_created_idx
ON analytics_events(vendor_id, event_type, created_at DESC);

CREATE INDEX IF NOT EXISTS analytics_events_product_created_idx
ON analytics_events(product_id, created_at DESC)
WHERE product_id IS NOT NULL;
