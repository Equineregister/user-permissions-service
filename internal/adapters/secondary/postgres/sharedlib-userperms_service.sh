function test_psql_exitcode {
    #$1 = PSQL_EXITCODE
    if [[ $1 -eq 3 ]]; then 
        echo "psql exit code: $1"
        echo "Error occurred in a script/SQL statement"
        echo "Check any SQL statments being run"
        exit 1
    fi
    if [[ $1 -eq 2 ]]; then 
        echo "psql exit code: $1"
        echo "Connection to the server went bad" #from man psq
        echo "Check user/password/hostname/databasename"
        exit 1
    fi
    if [[ $1 -eq 1 ]]; then 
        echo "psql exit code: $1"
        echo "A fatal error in psql (e.g., out of memory, file not found)" #from man psql
        exit 1
    fi
    #TODO: Need to Exit if code != 0 as well - then just put unknown error 
}

