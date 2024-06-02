CREATE TABLE "waters" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "volume" smallint,
  "drank_at" timestamp DEFAULT current_timestamp
)
