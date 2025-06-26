CREATE TABLE server_settings (
    id SERIAL PRIMARY KEY,
    server_name VARCHAR(21) NOT NULL,
    server_id SMALLINT UNIQUE NOT NULL,
    settings JSONB,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_server_settings_server_name ON server_settings(server_name);
CREATE INDEX idx_server_settings_server_id ON server_settings(server_id);
CREATE INDEX idx_server_settings_created_at ON server_settings(created_at);
CREATE INDEX idx_server_settings_updated_at ON server_settings(updated_at);
CREATE INDEX idx_server_settings_deleted_at ON server_settings(deleted_at);

CREATE INDEX idx_server_settings_settings_gin ON server_settings USING GIN (settings);
