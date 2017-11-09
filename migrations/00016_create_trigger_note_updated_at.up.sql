CREATE TRIGGER update_note_updated_at
BEFORE UPDATE ON note
FOR EACH ROW
  BEGIN
    SELECT current_timestamp
    INTO :new.updated_at
    FROM dual;
  END;