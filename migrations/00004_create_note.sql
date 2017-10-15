CREATE TABLE note (
  id         NUMBER PRIMARY KEY,
  title      VARCHAR(100) NOT NULL,
  body       BLOB         NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP DEFAULT current_timestamp
)
