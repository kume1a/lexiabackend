-- Update existing folder types from old enum values to new ones
-- This is a manual migration to handle the enum value change from CONTAINER to FOLDER_COLLECTION

-- Update any existing folders with type 'CONTAINER' to 'FOLDER_COLLECTION'
UPDATE "folders" SET "type" = 'FOLDER_COLLECTION' WHERE "type" = 'CONTAINER';

-- Update the default value for the type column
ALTER TABLE "folders" ALTER COLUMN "type" SET DEFAULT 'WORD_COLLECTION';
