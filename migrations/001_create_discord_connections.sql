BEGIN;

CREATE TABLE discord_connections (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    webhook_url_encrypted BYTEA NOT NULL,

    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notify_todo_created BOOLEAN NOT NULL DEFAULT TRUE,
    notify_todo_completed BOOLEAN NOT NULL DEFAULT TRUE,
    notify_before_due BOOLEAN NOT NULL DEFAULT TRUE,

    remind_before_minutes INTEGER NOT NULL DEFAULT 60,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT discord_connections_user_unique
        UNIQUE (user_id),

    CONSTRAINT discord_connections_remind_before_check
        CHECK (remind_before_minutes BETWEEN 1 AND 10080)
);

COMMIT;
