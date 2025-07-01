-- Create "folder_subfolders" table
CREATE TABLE "folder_subfolders" (
  "folder_id" uuid NOT NULL,
  "parent_id" uuid NOT NULL,
  PRIMARY KEY ("folder_id", "parent_id"),
  CONSTRAINT "folder_subfolders_folder_id" FOREIGN KEY ("folder_id") REFERENCES "folders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "folder_subfolders_parent_id" FOREIGN KEY ("parent_id") REFERENCES "folders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
