CREATE OR REPLACE FUNCTION getNUmOFProjects
    RETURN NUMBER
IS
    l_ret NUMBER := 0;
    BEGIN
        SELECT count(*)
        INTO l_ret
        FROM project;
        RETURN l_ret;
    END;
