CREATE TYPE character_status AS ENUM ('active', 'locked', 'deleted');

CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE, 
    character_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    
    name VARCHAR(21) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    class INTEGER NOT NULL DEFAULT 0,
    experience_points BIGINT NOT NULL DEFAULT 0,
    woonz BIGINT NOT NULL DEFAULT 0,
    character_data JSONB,
    rebirth INTEGER NOT NULL DEFAULT 0,
    
    status character_status NOT NULL DEFAULT 'active',
    is_online BOOLEAN DEFAULT false,
    last_login TIMESTAMP WITH TIME ZONE,
    last_logout TIMESTAMP WITH TIME ZONE,
    
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_name CHECK (name ~ '^[a-zA-Z0-9_-]{3,21}$')
);

CREATE INDEX idx_characters_account_id ON characters(account_id);
CREATE INDEX idx_characters_name ON characters(name);
CREATE INDEX idx_characters_status ON characters(status);
CREATE INDEX idx_characters_is_online ON characters(is_online);
CREATE INDEX idx_characters_level ON characters(level);
CREATE INDEX idx_characters_class ON characters(class);
CREATE INDEX idx_characters_experience_points ON characters(experience_points);
CREATE INDEX idx_characters_rebirth ON characters(rebirth);
CREATE INDEX idx_characters_last_login ON characters(last_login);
CREATE INDEX idx_characters_last_logout ON characters(last_logout);
CREATE INDEX idx_characters_created_at ON characters(created_at);
CREATE INDEX idx_characters_updated_at ON characters(updated_at);

CREATE INDEX idx_characters_account_status ON characters(account_id, status);
CREATE INDEX idx_characters_online_status ON characters(is_online, status);
CREATE INDEX idx_characters_level_class ON characters(level, class);
CREATE INDEX idx_characters_exp_level ON characters(experience_points, level);

CREATE INDEX idx_characters_active_only ON characters(account_id, name, level, class) 
    WHERE status = 'active';

CREATE INDEX idx_characters_online_only ON characters(account_id, name, level, class) 
    WHERE is_online = true AND status = 'active';

CREATE INDEX idx_characters_character_data_gin ON characters USING GIN (character_data);
