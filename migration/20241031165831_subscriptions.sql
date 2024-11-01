-- migrate:up
  CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(26) NOT NULL,
    user_id VARCHAR(26) not null,
    subscription_plan_id int not null,
    stripe_subscription_id VARCHAR(255) not null,

    started_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    canceled_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT subscriptions_id_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id)
  );

-- migrate:down
  DROP TABLE IF EXISTS subscriptions;