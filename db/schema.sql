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

CREATE TABLE
    "decks" (
        "id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        "name" VARCHAR(255) NOT NULL,
        "description" TEXT DEFAULT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "updated_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "deleted_at" TIMESTAMP DEFAULT NULL,
        "user_id" integer NOT NULL,
        "public" BOOLEAN DEFAULT FALSE
    );