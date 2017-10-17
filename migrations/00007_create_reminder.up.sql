CREATE TABLE reminder (
  id         NUMBER PRIMARY KEY,
  title      VARCHAR(1000) NOT NULL,
  display_at TIMESTAMP     NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP DEFAULT current_timestamp
)
