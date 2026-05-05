CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS uuid UUID;

UPDATE users
SET uuid = gen_random_uuid()
WHERE uuid IS NULL;

ALTER TABLE users
ALTER COLUMN uuid SET DEFAULT gen_random_uuid();

ALTER TABLE users
ALTER COLUMN uuid SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS users_uuid_idx ON users(uuid);

CREATE TABLE IF NOT EXISTS role_upgrade_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
  requested_role_id INT NOT NULL REFERENCES roles(id),
  status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending | approved | rejected
  reviewed_by UUID REFERENCES users(uuid),
  reviewed_at TIMESTAMPTZ,
  created_at timestamp(0) with time zone NOT NULL DEFAULt now(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULt now()
);