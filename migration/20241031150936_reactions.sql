-- migrate:up
  CREATE TABLE IF NOT EXISTS reactions (
    id VARCHAR(26) NOT NULL,
    user_id VARCHAR(26),
    matched_user_id VARCHAR(26),
    type VARCHAR(35),
    matched_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT reactions_id_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (matched_user_id) REFERENCES users(id)
  );

-- migrate:down
  DROP TABLE IF EXISTS reactions;