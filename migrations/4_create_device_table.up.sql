CREATE TABLE device (
  user_id           INT NOT NULL,
  device_identifier VARCHAR(100) UNIQUE,
  FOREIGN KEY (user_id) REFERENCES user_profile (id),
  PRIMARY KEY (user_id, device_identifier)
);

