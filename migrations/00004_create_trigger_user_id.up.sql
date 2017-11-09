CREATE OR REPLACE TRIGGER user_profile_id_trg
BEFORE INSERT
  ON user_profile
FOR EACH ROW

  BEGIN
    SELECT user_profile_id_seq.NEXTVAL
    INTO :new.id
    FROM dual;
  END;