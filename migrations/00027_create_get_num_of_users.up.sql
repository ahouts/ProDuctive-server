CREATE OR REPLACE FUNCTION getNumOfUsers
    RETURN NUMBER
IS
    l_ret NUMBER := 0;
    BEGIN
        SELECT count(*)
        INTO l_ret
        FROM user_profile;
        RETURN l_ret;
    END;
