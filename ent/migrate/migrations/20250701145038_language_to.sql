-- Modify "folders" table
ALTER TABLE "folders" ALTER COLUMN "type" SET DEFAULT 'WORD_COLLECTION', ADD COLUMN "language_to" character varying NULL;
