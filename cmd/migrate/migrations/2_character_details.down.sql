DROP INDEX IF EXISTS idx_characters_character_data_gin;
DROP INDEX IF EXISTS idx_characters_online_only;
DROP INDEX IF EXISTS idx_characters_active_only;
DROP INDEX IF EXISTS idx_characters_exp_level;
DROP INDEX IF EXISTS idx_characters_level_class;
DROP INDEX IF EXISTS idx_characters_online_status;
DROP INDEX IF EXISTS idx_characters_account_last_used;
DROP INDEX IF EXISTS idx_characters_account_status;
DROP INDEX IF EXISTS idx_characters_updated_at;
DROP INDEX IF EXISTS idx_characters_created_at;
DROP INDEX IF EXISTS idx_characters_last_logout;
DROP INDEX IF EXISTS idx_characters_last_login;
DROP INDEX IF EXISTS idx_characters_rebirth;
DROP INDEX IF EXISTS idx_characters_experience_points;
DROP INDEX IF EXISTS idx_characters_class;
DROP INDEX IF EXISTS idx_characters_level;
DROP INDEX IF EXISTS idx_characters_is_last_used;
DROP INDEX IF EXISTS idx_characters_is_online;
DROP INDEX IF EXISTS idx_characters_status;
DROP INDEX IF EXISTS idx_characters_name;
DROP INDEX IF EXISTS idx_characters_account_id;

DROP TABLE IF EXISTS characters;

DROP TYPE IF EXISTS character_status;
