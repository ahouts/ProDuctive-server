CREATE OR REPLACE FUNCTION
  get_notes_for_project(project_id IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR nid IN (SELECT id
                FROM note
                WHERE note.project_id = project_id)
    LOOP
      PIPE ROW (nid.id);
    END LOOP;
    RETURN;
  END;
