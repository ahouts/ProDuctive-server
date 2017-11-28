CREATE OR REPLACE FUNCTION
  user_has_permission_for_note(m IN NUMBER, n IN number)
  RETURN array
  PIPELINED
AS
  cursor users_for_note is (select * from table(get_users_for_note(n)));
  cursor projects_for_note is (select * from table(get_users_for_project(get_project_for_note(n))));
  BEGIN
    FOR uid IN users_for_note
    LOOP
      if uid.column_value = m THEN
        PIPE ROW(1);
        RETURN;
      END IF;
    END LOOP;
    FOR uid IN projects_for_note
    LOOP
      if uid.column_value = m THEN
        PIPE ROW(1);
        RETURN;
      END IF;
    END LOOP;
    PIPE ROW(0);
    RETURN;
  END;
