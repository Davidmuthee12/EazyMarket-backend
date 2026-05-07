CREATE TABLE IF NOT EXISTS vendor_profiles (
  id              bigserial,
  user_id         UUID UNIQUE NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
  store_name      VARCHAR(150) NOT NULL,
  subdomain       VARCHAR(100) UNIQUE NOT NULL,
  description     TEXT,
  logo_url        TEXT,
  banner_url      TEXT,
  business_email  VARCHAR(255),
  business_phone  VARCHAR(20),
  address         TEXT,
  status          VARCHAR(20) DEFAULT 'pending', -- pending | approved | suspended
  commission_rate NUMERIC(5,2) DEFAULT 10.00,    
  created_at      TIMESTAMPTZ DEFAULT NOW(),
  updated_at      TIMESTAMPTZ DEFAULT NOW()
);