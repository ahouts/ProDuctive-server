CREATE OR REPLACE FUNCTION getAvgNotePerProject
    RETURN VARCHAR
IS
    l_temp       NUMBER := 0;
    l_numOfUsers NUMBER := 0;
    l_retStr     VARCHAR(2000) := 'The average number of notes per user is: ';
    e            CHAR := chr(10);
    BEGIN
        l_numOfUsers := getNumOfUsers();
        l_temp := getNumOfNotes() / l_numOfUsers;
        l_retStr := l_retStr || ROUND(l_temp, 2) || e;
        l_retStr := l_retStr || ' The average number of projects per user is: ';
        l_temp := getNumOfProjects() / l_numOfUsers;
        l_retStr := l_retStr || ROUND(l_temp, 2) || e;
        l_retStr := l_retStr || ' The average number of reminders per user is ';
        l_temp := getNumOfReminders() / l_numOfUsers;
        l_retStr := l_retStr || ROUND(l_temp, 2) || e;
        RETURN l_retStr;
    END;
