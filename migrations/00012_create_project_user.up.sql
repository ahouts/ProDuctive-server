CREATE TABLE project_user (
  user_id          NUMBER,
  project_id       NUMBER,
  permission_level NUMBER NOT NULL,
  CONSTRAINT pu_user_id_fk FOREIGN KEY (user_id) REFERENCES user_profile (id) ON DELETE CASCADE,
  CONSTRAINT pu_project_id_fk FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, project_id)
)