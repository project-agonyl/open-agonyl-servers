DROP TRIGGER IF EXISTS update_web_sessions_updated_at ON web_sessions;

DROP INDEX IF EXISTS idx_web_sessions_expires_at;
DROP INDEX IF EXISTS idx_web_sessions_session_id;
DROP INDEX IF EXISTS idx_web_sessions_account_id;

DROP TABLE IF EXISTS web_sessions; 