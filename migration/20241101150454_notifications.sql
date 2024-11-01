-- migrate:up
  CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(128) NOT NULL,
    user_id VARCHAR(26) not null,
    content JSON,
    is_read boolean,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT notifications_id_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
  );

-- migrate:down
  DROP TABLE IF EXISTS notifications;