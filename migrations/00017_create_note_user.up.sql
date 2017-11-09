CREATE TABLE note_user (
user_id          NUMBER,
note_id          NUMBER,
permission_level NUMBER NOT NULL,
CONSTRAINT n_user_id_fk FOREIGN KEY (user_id) REFERENCES user_profile (id) ON DELETE CASCADE,
CONSTRAINT n_note_id_fk FOREIGN KEY (note_id) REFERENCES note (id) ON DELETE CASCADE,
PRIMARY KEY (user_id, note_id)
)