CREATE OR REPLACE TRIGGER project_id_trg
BEFORE INSERT
  ON project
FOR EACH ROW

  BEGIN
    SELECT project_id_seq.NEXTVAL
    INTO :new.id
    FROM dual;
  END;