CREATE TABLE project (
  id         NUMBER PRIMARY KEY,
  title      VARCHAR(400) NOT NULL,
  owner_id   NUMBER       NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP DEFAULT current_timestamp,
  CONSTRAINT p_owner_id_fk FOREIGN KEY (owner_id) REFERENCES user_profile (id) ON DELETE CASCADE
)