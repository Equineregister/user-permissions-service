#!/bin/bash

# Script to run revert migration of a Tenant's userperms_service database schema (DDL)  NOT data changes

source ./sharedlib-userperms_service.sh

echo ""

function show_help {
    echo    ""
    echo -n "schemarevert-userperms_service.sh -h|--hostname <hostname> -u|--username <username> "
    echo -n "-d|--database <database> -r|--revertfile <revertfile> "
    echo    ""
    echo    ""
    echo    "All options are required"
    echo    "Note: PGPASSWORD is expected to be set as an env variable"
}


#POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--hostname)
      DBHOSTNAME="$2"
      shift # past argument
      shift # past value
      ;;
    -u|--username)
      USERNAME="$2"
      shift # past argument
      shift # past value
      ;;
    -d|--database)
      DATABASE="$2"
      shift # past argument
      shift # past value
      ;;
    -r|--revertfile)
      REVERTFILE="$2"
      shift # past argument
      shift # past value
      ;;
    -*|--*)
      echo "Unknown option \"$1\" exiting"
      show_help
      exit 1
      ;;
    *)
      #POSITIONAL_ARGS+=("$1") # save positional arg
      #shift # past argument
      echo "Extra options found \"$1\" exiting"
      show_help
      exit 1
      ;;
  esac
done

#set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

if [ -z "${DBHOSTNAME}" ]; then
    echo "DBHOSTNAME not found"
    show_help
    exit 1
fi
if [ -z "${USERNAME}" ]; then
    echo "USERNAME not found"
    show_help
    exit 1
fi
if [ -z "${DATABASE}" ]; then
    echo "DATABASE not found"
    show_help
    exit 1
fi

if [ -z "${REVERTFILE}" ]; then
    echo "REVERTFILE not set"
    show_help
    exit 1
fi
if [ ! -f "reverts/${REVERTFILE}" ]; then
    echo "Cannot find REVERTFILE: $REVERTFILE"
    exit 1
fi

if [[ $(head -1 "reverts/$REVERTFILE") == *$'\r' ]]; then 
    echo "Revertfile file contains CR please convert to UNIX line endings"
    exit 1
fi 


if [[ ${REVERTFILE:4:8} != "-revert-" ]]; then 
    echo "Invalid revertfile filename"
    exit 1
fi

REVERTFILE_NUMBERSTRING=${REVERTFILE:0:4}
REVERTFILE_NUMBER=$(echo $REVERTFILE_NUMBERSTRING | sed 's/^0*//')
if ! [[ $REVERTFILE_NUMBER =~ ^[0-9]+$ ]]; then
    echo "REVERTFILE_NUMBER not a number"
    exit 1
fi


if [ -z "${PGPASSWORD}" ]; then
    echo "PGPASSWORD env variable is not set"
    show_help
    exit 1
fi

PSQL=`which psql`
if [ ! -x "${PSQL}" ]; then
    echo "Unable to find psql"
    echo ""
    echo "Please install before running"
    exit 1
fi


CURRENTDB_NUMBER=`echo "SELECT current_schema_number FROM public.db_schema_number;" \
                   | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?



CURRENTDB_NUMBER=$(echo $CURRENTDB_NUMBER | sed 's/^0*//')
if ! [[ $CURRENTDB_NUMBER =~ ^[0-9]+$ ]]; then
    echo "CURRENTDB_NUMBER not a number"
    exit 1
fi

if [[ $REVERTFILE_NUMBER -ne $CURRENTDB_NUMBER ]]; then 
    echo "REVERTFILE_NUMBER: $REVERTFILE_NUMBER is not the same as the CURRENTDB_NUMBER: $CURRENTDB_NUMBER"
    exit 1
fi

#need to deal with migration no 1 which creates the db_schema_number table and inserts the current_schema_number
#if [[ $REVERTFILE_NUMBER -eq 1 ]]; then 
   # echo "This is the first migration revert please run below manually:"
   # cat reverts/$REVERTFILE
   # exit 0
##fi


let "REVERTFILE_NUMBER_MINUS_ONE=REVERTFILE_NUMBER-1"

##Pad REVERTFILE_NUMBER_MINUS_ONE with leading 0s 
if [[ $REVERTFILE_NUMBER_MINUS_ONE -lt 10 ]]; then 
    REVERTFILE_NUMBER_MINUS_ONE="000$REVERTFILE_NUMBER_MINUS_ONE"
elif [[ $REVERTFILE_NUMBER_MINUS_ONE -lt 100 ]]; then 
    REVERTFILE_NUMBER_MINUS_ONE="00$REVERTFILE_NUMBER_MINUS_ONE"
elif [[ $REVERTFILE_NUMBER_MINUS_ONE -lt 1000 ]]; then 
    REVERTFILE_NUMBER_MINUS_ONE="0$REVERTFILE_NUMBER_MINUS_ONE"
fi

REVERTFILE_SQL=`cat reverts/$REVERTFILE`

if [[ $REVERTFILE_NUMBER -eq 1 ]]; then 
REVERT_SQL=$(cat << EOF
BEGIN;
$REVERTFILE_SQL
COMMIT;
EOF
)
else
REVERT_SQL=$(cat << EOF
BEGIN;
$REVERTFILE_SQL
UPDATE public.db_schema_number SET current_schema_number='$REVERTFILE_NUMBER_MINUS_ONE';
INSERT INTO db_changes_log (filename, message, deploy_start) VALUES ('$REVERTFILE', '$REVERTFILE_NUMBERSTRING reverted to $REVERTFILE_NUMBER_MINUS_ONE',DEPLOY_START_TIMESTAMP);
COMMIT;
EOF
)
fi 


if [[ $TESTING -ne 1 ]]; then
#show below to user and get do confirmation prompt so they see SQL about to be run and machine etc.
    echo "HOSTNAME = $DBHOSTNAME" 
    echo "USERNAME = $USERNAME"
    echo "DATABASE = $DATABASE"
    echo "REVERTFILE = $REVERTFILE"
    echo "REVERTFILE_NUMBER = $REVERTFILE_NUMBER"
    echo "CURRENTDB_NUMBER = $CURRENTDB_NUMBER"
    echo ""
    echo "REVERT_SQL: "
    echo "$REVERT_SQL"
    echo ""
fi
if [ -z "${SKIPCONFIRM}" ]; then
    echo "Please confirm revert?"
    read confirm_revert
    if [ "$confirm_revert" != "yes" ]; then 
        echo "Revert not confirmed.  Exiting"
        exit
    fi
fi




DEPLOY_START=`date +%s`
if [[ $TESTING -ne 1 ]]; then
  echo $DEPLOY_START
fi


if [[ $REVERTFILE_NUMBER -eq 1 ]]; then 
REVERT_SQL=$(cat << EOF
BEGIN;
$REVERTFILE_SQL
COMMIT;
EOF
)
else
REVERT_SQL=$(cat << EOF
BEGIN;
$REVERTFILE_SQL
UPDATE public.db_schema_number SET current_schema_number='$REVERTFILE_NUMBER_MINUS_ONE';
INSERT INTO db_changes_log (filename, message, deploy_start) VALUES ('$REVERTFILE', '$REVERTFILE_NUMBERSTRING reverted to $REVERTFILE_NUMBER_MINUS_ONE',to_timestamp($DEPLOY_START) at time zone 'utc');
COMMIT;
EOF
)
fi 



NEW_REVERTDB_NUMBER=`echo "$REVERT_SQL" \
                        | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?

DEPLOY_END=`date +%s`
if [[ $TESTING -ne 1 ]]; then
  echo $DEPLOY_END
fi

if [[ $REVERTFILE_NUMBER -ne 1 ]]; then 

  #confirm current_schema_number is now correct
  NEW_REVERTDB_NUMBER=`echo "SELECT current_schema_number FROM public.db_schema_number;" \
                       | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
  test_psql_exitcode $?


  NEW_REVERTDB_NUMBER=$(echo $NEW_REVERTDB_NUMBER | sed 's/^0*//')
  if ! [[ $NEW_REVERTDB_NUMBER =~ ^[0-9]+$ ]]; then
    echo "ERROR - NEW_REVERTDB_NUMBER: $NEW_REVERTDB_NUMBER not a number"
    echo "Manually check on what has gone wrong."
    #do we log this in history table?
    exit 1
  fi

  let "REVERTFILE_NUMBER_CHECK=REVERTFILE_NUMBER-1"
  if [[ $REVERTFILE_NUMBER_CHECK -ne $NEW_REVERTDB_NUMBER ]]; then 
      echo "Revert applied but NEW_REVERTDB_NUMBER: $NEW_REVERTDB_NUMBER is not one less than REVERTFILE_NUMBER: $REVERTFILE_NUMBER"
      echo "Manually check on what has gone wrong."
      #do we log this in history table?
      exit 1
  fi
  RES=`echo "UPDATE db_changes_log SET deploy_end=to_timestamp($DEPLOY_END) at time zone 'utc' WHERE filename='$REVERTFILE' AND deploy_start=to_timestamp($DEPLOY_START) at time zone 'utc';" \
       | psql -v ON_ERROR_STOP=1 -1qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
  test_psql_exitcode $?
  echo "$RES"

fi


echo "Revert from $REVERTFILE_NUMBER to $NEW_REVERTDB_NUMBER has been applied successfully"
if [[ $TESTING -ne 1 ]]; then
    echo " to $DATABASE on $HOSTNAME using $REVERTFILE"
fi
