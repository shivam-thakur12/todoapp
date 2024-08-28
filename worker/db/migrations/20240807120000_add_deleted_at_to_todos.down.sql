BEGIN;
-- Remove the deleted_at column from the todos table
ALTER TABLE todos
DROP COLUMN deleted_at;
COMMIT;
