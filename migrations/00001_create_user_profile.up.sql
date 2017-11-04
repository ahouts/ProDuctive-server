CREATE TABLE user_profile (
  id            NUMBER PRIMARY KEY,
  email         VARCHAR(100) UNIQUE,
  password_hash RAW(60) NOT NULL,
  created_at    TIMESTAMP DEFAULT current_timestamp,
  updated_at    TIMESTAMP DEFAULT current_timestamp
)
