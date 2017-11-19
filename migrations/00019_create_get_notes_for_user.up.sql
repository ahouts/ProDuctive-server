CREATE OR REPLACE FUNCTION
  get_notes_for_user(user_id IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR nid IN (SELECT id
                FROM note
                WHERE note.owner_id = user_id)
    LOOP
      PIPE ROW (nid.id);
    END LOOP;
    FOR nid IN (SELECT note_id
                FROM note_user
                WHERE note_user.user_id = user_id)
    LOOP
      PIPE ROW (nid.note_id);
    END LOOP;
    RETURN;
  END;
