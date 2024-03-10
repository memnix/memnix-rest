-- Create "users" table
CREATE TABLE "users" ("id" integer NOT NULL GENERATED ALWAYS AS IDENTITY, "email" character varying(255) NOT NULL, "password" character varying(100) NOT NULL, "created_at" timestamp NOT NULL DEFAULT (now() AT TIME ZONE 'utc'::text), "updated_at" timestamp NOT NULL DEFAULT (now() AT TIME ZONE 'utc'::text), "deleted_at" timestamp NULL, "username" character varying(50) NULL DEFAULT NULL::character varying, "oauth_id" character varying(255) NULL DEFAULT NULL::character varying, "oauth_provider" character varying(50) NULL DEFAULT NULL::character varying, "has_oauth" boolean NULL DEFAULT false, "avatar" character varying(255) NULL DEFAULT NULL::character varying, "permission" integer NULL DEFAULT 0, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create index "users_oauth_id_key" to table: "users"
CREATE UNIQUE INDEX "users_oauth_id_key" ON "users" ("oauth_id");
