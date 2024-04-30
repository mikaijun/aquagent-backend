CREATE TABLE "question" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "title" varchar,
  "content" varchar,
  "file_path" varchar,
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp DEFAULT current_timestamp,
  "deleted_at" timestamp
)
