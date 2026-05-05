ALTER TABLE user_invitation
DROP CONSTRAINT IF EXISTS user_invitation_user_uuid_fkey;

ALTER TABLE user_invitation
ALTER COLUMN token TYPE BYTEA USING decode(token, 'hex');

ALTER TABLE user_invitation
DROP COLUMN IF EXISTS expiry;

ALTER TABLE user_invitation
DROP COLUMN IF EXISTS user_uuid;
