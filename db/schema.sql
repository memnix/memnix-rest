CREATE TABLE
    "users" (
        "id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        "email" VARCHAR(255) NOT NULL UNIQUE,
        "password" VARCHAR(100) NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "updated_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "deleted_at" TIMESTAMP DEFAULT NULL,
        "username" VARCHAR(50) DEFAULT NULL,
        "oauth_id" VARCHAR(255) DEFAULT NULL UNIQUE,
        "oauth_provider" VARCHAR(50) DEFAULT NULL,
        "has_oauth" BOOLEAN DEFAULT FALSE,
        "avatar" VARCHAR(255) DEFAULT NULL,
        "permission" integer DEFAULT 0
    );