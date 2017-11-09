CREATE TRIGGER update_reminder_updated_at
BEFORE UPDATE ON reminder
FOR EACH ROW
  BEGIN
    SELECT current_timestamp
    INTO :new.updated_at
    FROM dual;
  END;