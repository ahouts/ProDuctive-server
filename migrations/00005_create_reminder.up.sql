CREATE TABLE reminder (
  id         NUMBER PRIMARY KEY,
  user_id    NUMBER   NOT NULL,
  body       LONG RAW NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP DEFAULT current_timestamp,
  CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES user_profile (id) ON DELETE CASCADE
)