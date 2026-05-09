-- PRODUCTS
CREATE TABLE IF NOT EXISTS products (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  vendor_id       UUID NOT NULL REFERENCES vendor_profiles(user_id) ON DELETE CASCADE,
  category_id     UUID REFERENCES categories(id),
  name            VARCHAR(255) NOT NULL,
  slug            VARCHAR(255) UNIQUE NOT NULL,
  description     TEXT,
  price           NUMERIC(12,2) NOT NULL,
  compare_price   NUMERIC(12,2),
  stock_quantity  INT NOT NULL DEFAULT 0,
  sku             VARCHAR(100),
  status          VARCHAR(20) DEFAULT 'draft',  -- draft | published | archived
  weight          NUMERIC(8,2),
  created_at      timestamp(0) with time zone NOT NULL DEFAULt now(),
  updated_at      timestamp(0) with time zone NOT NULL DEFAULt now()
);

CREATE TABLE IF NOT EXISTS product_images (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  url        TEXT NOT NULL,
  alt_text   TEXT,
  sort_order INT DEFAULT 0,
  is_primary BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_tags (
  product_id UUID REFERENCES products(id) ON DELETE CASCADE,
  tag        VARCHAR(50),
  PRIMARY KEY (product_id, tag)
);