DROP TRIGGER IF EXISTS create_profile_trigger ON accounts;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;

DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS generate_account_id();
DROP FUNCTION IF EXISTS create_profile_for_account();

DROP INDEX IF EXISTS idx_accounts_username;
DROP INDEX IF EXISTS idx_accounts_email;
DROP INDEX IF EXISTS idx_accounts_account_id;
DROP INDEX IF EXISTS idx_accounts_status;
DROP INDEX IF EXISTS idx_accounts_subscription_tier;
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_accounts_last_login;
DROP INDEX IF EXISTS idx_accounts_guild_id;

DROP INDEX IF EXISTS idx_profiles_account_id;
DROP INDEX IF EXISTS idx_profiles_display_name;
DROP INDEX IF EXISTS idx_profiles_country;
DROP INDEX IF EXISTS idx_profiles_is_public;

DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS accounts;

DROP TYPE IF EXISTS account_status;
