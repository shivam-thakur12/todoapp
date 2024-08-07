-- Add the deleted_at column to todos table
ALTER TABLE todos
ADD COLUMN deleted_at TIMESTAMP;
