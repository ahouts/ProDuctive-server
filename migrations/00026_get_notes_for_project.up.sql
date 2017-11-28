CREATE OR REPLACE FUNCTION
  get_notes_for_project(m IN NUMBER)
  RETURN ARRAY
PIPELINED
AS
  BEGIN
    FOR nid IN (SELECT id
                FROM note
                WHERE note.project_id = m)
    LOOP
      PIPE ROW (nid.id);
    END LOOP;
    RETURN;
  END;
