CREATE OR REPLACE FUNCTION
  permission_for_project(user_id IN NUMBER, project_id IN number)
  RETURN array
PIPELINED
AS
  cursor users_for_project is (select * from table(get_users_for_project(project_id)));
  BEGIN
    FOR uid IN users_for_project
    LOOP
      if uid.column_value = user_id THEN
        PIPE ROW(1);
        RETURN;
      END IF;
    END LOOP;
    PIPE ROW(0);
    RETURN;
  END;
