CREATE TABLE
    "users" (
        "id" SERIAL PRIMARY KEY,
        "email" VARCHAR(255) NOT NULL UNIQUE,
        "password" VARCHAR(255) NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "updated_at" TIMESTAMP NOT NULL DEFAULT (now () at time zone 'utc'),
        "deleted_at" TIMESTAMP DEFAULT NULL,
        "username" VARCHAR(255) DEFAULT NULL
    );