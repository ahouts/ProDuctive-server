CREATE OR REPLACE FUNCTION getNumOfReminders
    RETURN NUMBER
IS
    l_ret NUMBER := 0;
    BEGIN
        SELECT count(*)
        INTO l_ret
        FROM reminder;
        RETURN l_ret;
    END;
