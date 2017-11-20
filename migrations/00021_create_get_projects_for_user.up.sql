CREATE OR REPLACE FUNCTION
  get_projects_for_user(user_id IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR nid IN (SELECT id
                FROM project
                WHERE project.owner_id = user_id)
    LOOP
      PIPE ROW (nid.id);
    END LOOP;
    FOR nid IN (SELECT project_id
                FROM project_user
                WHERE project_user.user_id = user_id)
    LOOP
      PIPE ROW (nid.project_id);
    END LOOP;
    RETURN;
  END;
