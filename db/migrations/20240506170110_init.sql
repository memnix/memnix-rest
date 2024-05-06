-- Modify "users" table
ALTER TABLE "users" ADD CONSTRAINT "users_email_key" UNIQUE USING INDEX "users_email_key", ADD CONSTRAINT "users_oauth_id_key" UNIQUE USING INDEX "users_oauth_id_key";
