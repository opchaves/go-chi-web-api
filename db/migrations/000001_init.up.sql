BEGIN;

SET timezone TO 'GMT';

CREATE TABLE IF NOT EXISTS users(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "verified" BOOLEAN NOT NULL DEFAULT false,
  "verification_token" VARCHAR,
  "avatar" VARCHAR NOT NULL DEFAULT 'user_avatar0.png',
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "uq_users_email" UNIQUE ("email"),
  CONSTRAINT "pk_users_id" PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS workspaces(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "currency" VARCHAR NOT NULL,
  "language" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_workspaces_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_workspaces_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS workspaces_users(
  "workspace_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  "role" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_workspaces_users_id" PRIMARY KEY ("workspace_id", "user_id")
);

CREATE TABLE IF NOT EXISTS categories(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "icon" VARCHAR,
  "color" VARCHAR,
  "cat_type" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "workspace_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_categories_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_catetories_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES "workspaces"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS accounts(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "balance" DECIMAL(10, 2) NOT NULL DEFAULT 0,
  "user_id" UUID NOT NULL,
  "workspace_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_accounts_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_accounts_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES "workspaces"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS tokens(
  "id" INT GENERATED ALWAYS AS IDENTITY,
  "token" VARCHAR NOT NULL,
  "identifier" VARCHAR,
  "mobile" BOOLEAN NOT NULL DEFAULT false,
  "user_id" UUID NOT NULL,
  "expires_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_tokens_id" PRIMARY KEY ("id"),
  CONSTRAINT "uq_tokens_token" UNIQUE ("token"),
  CONSTRAINT "fk_tokens_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

COMMIT;
