CREATE TABLE note (
  id         NUMBER PRIMARY KEY,
  title      VARCHAR(400)   NOT NULL,
  body       VARCHAR2(4000) NOT NULL,
  owner_id   NUMBER         NOT NULL,
  project_id NUMBER,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP DEFAULT current_timestamp,
  CONSTRAINT n_owner_id_fk FOREIGN KEY (owner_id) REFERENCES user_profile (id) ON DELETE CASCADE,
  CONSTRAINT n_project_id_fk FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE
)