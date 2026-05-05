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

ALTER TABLE user_invitation
ADD COLUMN IF NOT EXISTS user_uuid UUID;

ALTER TABLE user_invitation
ADD COLUMN IF NOT EXISTS expiry TIMESTAMPTZ;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'user_invitation'
          AND column_name = 'token'
          AND data_type = 'bytea'
    ) THEN
        ALTER TABLE user_invitation
        ALTER COLUMN token TYPE TEXT USING encode(token, 'hex');
    END IF;
END $$;

UPDATE user_invitation ui
SET user_uuid = u.uuid
FROM users u
WHERE ui.user_uuid IS NULL
  AND ui.user_id = u.id;

UPDATE user_invitation
SET expiry = NOW() + INTERVAL '3 days'
WHERE expiry IS NULL;

ALTER TABLE user_invitation
ALTER COLUMN user_uuid SET NOT NULL;

ALTER TABLE user_invitation
ALTER COLUMN expiry SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'user_invitation_user_uuid_fkey'
    ) THEN
        ALTER TABLE user_invitation
        ADD CONSTRAINT user_invitation_user_uuid_fkey
        FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE;
    END IF;
END $$;
