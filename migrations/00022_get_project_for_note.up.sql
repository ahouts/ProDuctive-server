CREATE OR REPLACE FUNCTION
  get_project_for_note(note_id in number)
  RETURN number
AS
  project_id note.project_id%type;
  BEGIN
    select note.project_id
    into project_id
    from note
    where note.id = note_id;
    return project_id;
  END;
