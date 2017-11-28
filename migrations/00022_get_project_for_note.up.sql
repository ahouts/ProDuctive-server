CREATE OR REPLACE FUNCTION
  get_project_for_note(m in number)
  RETURN number
AS
  project_id note.project_id%type;
  BEGIN
    select note.project_id
    into project_id
    from note
    where note.id = m;
    return project_id;
  END;
