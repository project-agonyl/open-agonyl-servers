CREATE TABLE item_sequences (
    id SERIAL PRIMARY KEY,
    server_id TEXT NOT NULL UNIQUE,
    last_allocated BIGINT NOT NULL,
    batch_size INTEGER NOT NULL DEFAULT 1000,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_item_sequences_server_id ON item_sequences(server_id);

CREATE OR REPLACE FUNCTION allocate_sequence_batch(
    p_server_id TEXT,
    p_batch_size INTEGER DEFAULT 1000
) RETURNS TABLE(start_id BIGINT, end_id BIGINT) AS $$
DECLARE
    current_max BIGINT;
    new_max BIGINT;
BEGIN
    SELECT COALESCE(MAX(last_allocated), 0) INTO current_max 
    FROM item_sequences;
    
    new_max := current_max + p_batch_size;
    
    INSERT INTO item_sequences (server_id, last_allocated, batch_size)
    VALUES (p_server_id, new_max, p_batch_size)
    ON CONFLICT (server_id) DO UPDATE SET
        last_allocated = new_max,
        batch_size = p_batch_size,
        updated_at = NOW();
    
    RETURN QUERY SELECT current_max + 1, new_max;
END;
$$ LANGUAGE plpgsql;