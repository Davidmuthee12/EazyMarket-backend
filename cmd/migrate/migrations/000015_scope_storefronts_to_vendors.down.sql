DROP INDEX IF EXISTS products_vendor_slug_idx;

ALTER TABLE products
ADD CONSTRAINT products_slug_key UNIQUE (slug);

DROP INDEX IF EXISTS carts_one_active_per_user_vendor_idx;

ALTER TABLE carts
DROP COLUMN IF EXISTS vendor_id;

ALTER TABLE carts
ADD CONSTRAINT carts_user_id_key UNIQUE (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS carts_one_active_per_user_idx
ON carts(user_id)
WHERE status = 'active';
