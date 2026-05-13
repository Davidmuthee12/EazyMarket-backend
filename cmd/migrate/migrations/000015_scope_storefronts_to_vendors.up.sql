ALTER TABLE carts
DROP CONSTRAINT IF EXISTS carts_user_id_key;

DROP INDEX IF EXISTS carts_one_active_per_user_idx;

ALTER TABLE carts
ADD COLUMN IF NOT EXISTS vendor_id UUID REFERENCES vendor_profiles(user_id) ON DELETE CASCADE;

UPDATE carts c
SET vendor_id = cart_vendors.vendor_id
FROM (
  SELECT DISTINCT ON (ci.cart_id)
    ci.cart_id,
    p.vendor_id
  FROM cart_items ci
  JOIN products p ON p.id = ci.product_id
  WHERE NOT EXISTS (
    SELECT 1
    FROM cart_items ci2
    JOIN products p2 ON p2.id = ci2.product_id
    WHERE ci2.cart_id = ci.cart_id
      AND p2.vendor_id <> p.vendor_id
  )
  ORDER BY ci.cart_id, p.vendor_id
) AS cart_vendors
WHERE c.id = cart_vendors.cart_id
  AND c.vendor_id IS NULL;

DELETE FROM cart_items
WHERE cart_id IN (
  SELECT id FROM carts WHERE vendor_id IS NULL
);

DELETE FROM carts
WHERE vendor_id IS NULL;

ALTER TABLE carts
ALTER COLUMN vendor_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS carts_one_active_per_user_vendor_idx
ON carts(user_id, vendor_id)
WHERE status = 'active';

ALTER TABLE products
DROP CONSTRAINT IF EXISTS products_slug_key;

CREATE UNIQUE INDEX IF NOT EXISTS products_vendor_slug_idx
ON products(vendor_id, slug);
