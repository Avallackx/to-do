CREATE TABLE IF NOT EXISTS "tasks" (
   "id" BIGINT PRIMARY KEY,
   "title" TEXT NOT NULL,
   "todo" TEXT NOT NULL,
   "completed" BOOLEAN DEFAULT FALSE,
   "created_at" TIMESTAMP NOT NULL DEFAULT 'now()',
   "updated_at" TIMESTAMP NOT NULL DEFAULT 'now()',
   "deleted_at" TIMESTAMP
)