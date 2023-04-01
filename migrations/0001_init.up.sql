BEGIN;
CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "telegram_id" bigint NOT NULL UNIQUE,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "is_subscribed_to_broadcasting" boolean NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMIT;