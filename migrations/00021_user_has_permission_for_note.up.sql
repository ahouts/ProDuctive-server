CREATE OR REPLACE FUNCTION
  user_has_permission_for_note(user_id IN NUMBER, note_id IN number)
  RETURN array
  PIPELINED
AS
  cursor users_for_note is (select * from table(get_users_for_note(note_id)));
  BEGIN
    FOR uid IN users_for_note
    LOOP
      if uid.column_value = user_id THEN
        PIPE ROW(1);
        RETURN;
      END IF;
    END LOOP;
    PIPE ROW(0);
    RETURN;
  END;
