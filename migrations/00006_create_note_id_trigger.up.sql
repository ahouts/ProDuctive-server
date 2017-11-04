CREATE OR REPLACE TRIGGER note_id_trg
  BEFORE INSERT
  ON note
  FOR EACH ROW

  BEGIN
    SELECT note_id_seq.NEXTVAL
    INTO :new.id
    FROM dual;
  END;