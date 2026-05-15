UPDATE vendor_profiles
SET subdomain = LOWER(TRIM(subdomain))
WHERE subdomain <> LOWER(TRIM(subdomain));

CREATE UNIQUE INDEX IF NOT EXISTS role_upgrade_requests_one_pending_per_user_idx
ON role_upgrade_requests(user_id)
WHERE status = 'pending';
