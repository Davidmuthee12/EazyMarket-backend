ALTER TABLE
    carts
ADD COLUMN 
    status VARCHAR(20) NOT NULL DEFAULT 'active';

CREATE UNIQUE INDEX IF NOT EXISTS carts_one_active_per_user_idx
ON carts(user_id)
WHERE status = 'active';

ALTER TABLE 
    cart_items
ADD COLUMN
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
ADD COLUMN
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now();
