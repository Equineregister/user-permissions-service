#!/bin/bash

# Script to run Tenant specific data population for a Tenant's userperms_service database.


source ./sharedlib-userperms_service.sh


echo ""

function show_help {
    echo    ""
    echo -n "schemaupdate-userperms_service.sh -h|--hostname <hostname> -u|--username <username> "
    echo -n "-d|--database <database> -m|--migrationfile <tenant specific migration file> "
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
    -m|--migrationfile)
      MIGRATIONFILE="$2"
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

if [ -z "${MIGRATIONFILE}" ]; then
    echo "MIGRATIONFILE not set"
    show_help
    exit 1
fi


if [ ! -f "tenants/${MIGRATIONFILE}" ]; then
    echo "Cannot find MIGRATIONFILE: tenants/$MIGRATIONFILE"
    exit 1
fi


if [[ $(head -1 "tenants/$MIGRATIONFILE") == *$'\r' ]]; then 
    echo "Migration file: tenants/$MIGRATIONFILE contains CR please convert to UNIX line endings"
    exit 1
fi 

if [[ ${MIGRATIONFILE:4:1} != "-" ]]; then 
    echo "Invalid migration filename: $MIGRATIONFILE"
    exit 1
fi

MIGRATIONFILE_NUMBERSTRING=${MIGRATIONFILE:0:4}
MIGRATIONFILE_NUMBER=$(echo $MIGRATIONFILE_NUMBERSTRING | sed 's/^0*//')
if ! [[ $MIGRATIONFILE_NUMBER =~ ^[0-9]+$ ]]; then
    echo "MIGRATIONFILE_NUMBER not a number: $MIGRATIONFILE_NUMBER from $MIGRATIONFILE"
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

# Note: We expect the db_datapop_number table to already exist.
MIGRATIONDB_NUMBERSTRING=`echo "SELECT current_datapop_number FROM public.db_datapop_number;" \
                    | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?

MIGRATIONDB_NUMBER=$(echo $MIGRATIONDB_NUMBERSTRING | sed 's/^0*//')
if ! [[ $MIGRATIONDB_NUMBER =~ ^[0-9]+$ ]]; then
    echo "MIGRATIONDB_NUMBER not a number"
    exit 1
fi

let "MIGRATIONFILE_NUMBER_CHECK=MIGRATIONFILE_NUMBER-1"
if [[ $MIGRATIONFILE_NUMBER_CHECK -ne $MIGRATIONDB_NUMBER ]]; then 
    echo "MIGRATIONFILE_NUMBER: $MIGRATIONFILE_NUMBER is not 1 step ahead of MIGRATIONDB_NUMBER: $MIGRATIONDB_NUMBER"
    exit 1
fi

MIGRATIONFILE_SQL=`cat tenants/$MIGRATIONFILE`

MIGRATION_SQL=$(cat << EOF
BEGIN;
$MIGRATIONFILE_SQL
UPDATE public.db_datapop_number SET current_datapop_number='$MIGRATIONFILE_NUMBERSTRING';
INSERT INTO db_changes_log (filename, message, deploy_start) VALUES ('$MIGRATIONFILE', '$MIGRATIONDB_NUMBERSTRING migrated to $MIGRATIONFILE_NUMBERSTRING',DEPLOY_START_TIMESTAMP);
COMMIT;
EOF
)

if [[ $TESTING -ne 1 ]]; then

    #show below to user and get do confirmation prompt so they see SQL about to be run and machine etc.
    echo "HOSTNAME = $DBHOSTNAME" 
    echo "USERNAME = $USERNAME"
    echo "DATABASE = $DATABASE"
    echo "MIGRATIONFILE = $MIGRATIONFILE"
    echo "MIGRATIONFILE_NUMBER = $MIGRATIONFILE_NUMBER"
    echo "MIGRATIONDB_NUMBER = $MIGRATIONDB_NUMBER"
    echo ""
    echo "MIGRATION_SQL: "
    echo "$MIGRATION_SQL"
    echo ""
fi
if [ -z "${SKIPCONFIRM}" ]; then
    echo "Please confirm?"
    read confirm_migration
    if [ "$confirm_migration" != "yes" ]; then 
        echo "Not confirmed.  Exiting"
        exit
    fi
fi


DEPLOY_START=`date +%s`
if [[ $TESTING -ne 1 ]]; then
  echo $DEPLOY_START
fi

MIGRATIONFILE_SQL=`cat tenants/$MIGRATIONFILE`

MIGRATION_SQL=$(cat << EOF
BEGIN;
$MIGRATIONFILE_SQL
UPDATE public.db_datapop_number SET current_datapop_number='$MIGRATIONFILE_NUMBERSTRING';
INSERT INTO db_changes_log (filename, message, deploy_start) VALUES ('$MIGRATIONFILE', '$MIGRATIONDB_NUMBERSTRING migrated to $MIGRATIONFILE_NUMBERSTRING',to_timestamp($DEPLOY_START) at time zone 'utc');
COMMIT;
EOF
)

NEW_MIGRATIONDB_NUMBER=`echo "$MIGRATION_SQL" \
                        | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?

DEPLOY_END=`date +%s`
if [[ $TESTING -ne 1 ]]; then
  echo $DEPLOY_END
fi

#confirm current_datapop_number is now correct
NEW_MIGRATIONDB_NUMBER=`echo "SELECT current_datapop_number FROM public.db_datapop_number;" \
                        | psql -v ON_ERROR_STOP=1 -qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?

NEW_MIGRATIONDB_NUMBER=$(echo $NEW_MIGRATIONDB_NUMBER | sed 's/^0*//')
if ! [[ $NEW_MIGRATIONDB_NUMBER =~ ^[0-9]+$ ]]; then
    echo "ERROR - NEW_MIGRATIONDB_NUMBER: $NEW_MIGRATIONDB_NUMBER not a number"
    echo "Manually check on what has gone wrong."
    #TODO: do we log this in history table?
    exit 1
fi
if [[ $MIGRATIONFILE_NUMBER -ne $NEW_MIGRATIONDB_NUMBER ]]; then 
    echo "Migration applied but NEW_MIGRATIONDB_NUMBER: $NEW_MIGRATIONDB_NUMBER is not the same as MIGRATIONFILE_NUMBER: $MIGRATIONFILE_NUMBER"
    echo "Manually check on what has gone wrong."
    #TODO: do we log this in history table?
    exit 1
fi

RES=`echo "UPDATE db_changes_log SET deploy_end=to_timestamp($DEPLOY_END) at time zone 'utc' WHERE filename='$MIGRATIONFILE' AND deploy_start=to_timestamp($DEPLOY_START) at time zone 'utc';" \
     | psql -v ON_ERROR_STOP=1 -1qAt -h $DBHOSTNAME -U $USERNAME -d $DATABASE`
test_psql_exitcode $?
echo "$RES"

echo "Migration from $MIGRATIONDB_NUMBER to $NEW_MIGRATIONDB_NUMBER has been applied successfully" 
if [[ $TESTING -ne 1 ]]; then
    echo "to $DATABASE on $DBHOSTNAME using $MIGRATIONFILE"
fi


