CREATE OR REPLACE FUNCTION
  get_users_for_note(note_id IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR uid IN (SELECT owner_id
                FROM note
                WHERE note.id = note_id)
    LOOP
      PIPE ROW (uid.owner_id);
    END LOOP;
    FOR uid IN (SELECT user_id
                FROM note_user
                WHERE note_user.note_id = note_id)
    LOOP
      PIPE ROW (uid.user_id);
    END LOOP;

    FOR uid IN (SELECT user_id
                FROM project_user
                WHERE project_user.project_id = get_project_for_note(note_id))
    LOOP
      PIPE ROW (uid.user_id);
    END LOOP;

    RETURN;
  END;
