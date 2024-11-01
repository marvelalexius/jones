-- migrate:up
  CREATE TABLE IF NOT EXISTS subscription_plans (
    id int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    name VARCHAR(128),
    price decimal(10,0),
    features text[],
    stripe_price_id VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT subscription_plans_id_pkey PRIMARY KEY (id)
  );

-- migrate:down
  DROP TABLE IF EXISTS subscription_plans;