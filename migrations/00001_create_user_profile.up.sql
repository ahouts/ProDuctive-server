CREATE TABLE user_profile (
  id            NUMBER PRIMARY KEY,
  email         VARCHAR(100) UNIQUE,
  password_hash CHAR(32) NOT NULL,
  created_at    TIMESTAMP DEFAULT current_timestamp,
  updated_at    TIMESTAMP DEFAULT current_timestamp
)
