CREATE TABLE "waters" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "volume" smallint,
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp DEFAULT current_timestamp
)
