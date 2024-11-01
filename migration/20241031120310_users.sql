-- migrate:up
  CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    bio TEXT,
    gender VARCHAR(10),
    preference VARCHAR(10),
    age INT,
    stripe_customer_id VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    
    CONSTRAINT users_id_pkey PRIMARY KEY (id)
  );

-- migrate:down
  DROP TABLE IF EXISTS users;