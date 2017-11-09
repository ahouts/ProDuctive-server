CREATE OR REPLACE TRIGGER reminder_id_trg
BEFORE INSERT
  ON reminder
FOR EACH ROW

  BEGIN
    SELECT reminder_id_seq.NEXTVAL
    INTO :new.id
    FROM dual;
  END;