DROP INDEX IF EXISTS idx_server_settings_settings_gin;
DROP INDEX IF EXISTS idx_server_settings_deleted_at;
DROP INDEX IF EXISTS idx_server_settings_updated_at;
DROP INDEX IF EXISTS idx_server_settings_created_at;
DROP INDEX IF EXISTS idx_server_settings_server_id;
DROP INDEX IF EXISTS idx_server_settings_server_name;

DROP TABLE IF EXISTS server_settings;
