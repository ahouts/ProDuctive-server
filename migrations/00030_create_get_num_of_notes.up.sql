CREATE OR REPLACE FUNCTION getNumOfNotes
    RETURN NUMBER
IS
    l_ret NUMBER := 0;
    BEGIN
        SELECT count(*)
        INTO l_ret
        FROM note;
        RETURN l_ret;
    END;
