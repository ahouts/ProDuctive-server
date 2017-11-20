CREATE OR REPLACE FUNCTION
  get_users_for_project(project_id IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR uid IN (SELECT owner_id
                FROM project
                WHERE project.id = project_id)
    LOOP
      PIPE ROW (uid.owner_id);
    END LOOP;
    FOR uid IN (SELECT user_id
                FROM project_user
                WHERE project_user.project_id = project_id)
    LOOP
      PIPE ROW (uid.user_id);
    END LOOP;
    RETURN;
  END;
