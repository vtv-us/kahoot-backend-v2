CREATE TABLE "user" (
  "user_id" text PRIMARY KEY,
  "email" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "password" text NOT NULL,
  "verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "google_id" text,
  "facebook_id" text
);

CREATE TABLE "group" (
  "group_id" text PRIMARY KEY,
  "group_name" text NOT NULL,
  "created_by" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_group" (
  "user_id" text NOT NULL,
  "group_id" text NOT NULL,
  "role" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "group_id")
);

CREATE INDEX ON "user" using btree("user_id");

CREATE INDEX ON "user" using btree("email");

CREATE INDEX ON "group" using btree("group_id");

CREATE INDEX ON "group" using btree("created_by");

CREATE INDEX ON "user_group" using btree("user_id", "group_id");

ALTER TABLE "user_group" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "user_group" ADD FOREIGN KEY ("group_id") REFERENCES "group" ("group_id");