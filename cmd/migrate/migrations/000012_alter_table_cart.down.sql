DROP INDEX IF EXISTS carts_one_active_per_user_idx;

ALTER TABLE carts
DROP COLUMN IF EXISTS status;

ALTER TABLE cart_items
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;
