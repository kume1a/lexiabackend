-- Modify "folders" table
ALTER TABLE "folders" ALTER COLUMN "language_from" DROP NOT NULL, ADD COLUMN "type" character varying NOT NULL DEFAULT 'word_collection';
