CREATE TABLE IF NOT EXISTS tokens(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "token" VARCHAR NOT NULL,
  "identifier" VARCHAR,
  "mobile" BOOLEAN NOT NULL DEFAULT false,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "pk_tokens_id" PRIMARY KEY ("id"),
  CONSTRAINT "uq_tokens_token" UNIQUE ("token"),
  CONSTRAINT "fk_tokens_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);
