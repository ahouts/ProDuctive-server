CREATE TRIGGER update_user_profile_updated_at
BEFORE UPDATE ON user_profile
FOR EACH ROW
  BEGIN
    SELECT current_timestamp
    INTO :new.updated_at
    FROM dual;
  END;