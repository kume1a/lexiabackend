-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "create_time" timestamptz NOT NULL,
  "update_time" timestamptz NOT NULL,
  "username" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "user_email" to table: "users"
CREATE UNIQUE INDEX "user_email" ON "users" ("email");
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create "folders" table
CREATE TABLE "folders" (
  "id" uuid NOT NULL,
  "create_time" timestamptz NOT NULL,
  "update_time" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "word_count" integer NOT NULL,
  "language_from" character varying NOT NULL,
  "user_folders" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "folders_users_folders" FOREIGN KEY ("user_folders") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "words" table
CREATE TABLE "words" (
  "id" uuid NOT NULL,
  "create_time" timestamptz NOT NULL,
  "update_time" timestamptz NOT NULL,
  "text" character varying NOT NULL,
  "definition" character varying NOT NULL,
  "folder_words" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "words_folders_words" FOREIGN KEY ("folder_words") REFERENCES "folders" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
