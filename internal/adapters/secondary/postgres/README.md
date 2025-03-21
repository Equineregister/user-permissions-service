# user_permissions_service

**_The system doesn't support schema changes per-Tenant, all Tenants must use the exact same schema._**

## Checklist for adding new migration

1.  Create .sql files for migration - named correctly - add to /migrations/ folder.

2.  Create .sql files for revert - named correctly - add to /reverts/ folder.

Note: Once commited and migration applied do not edit or update it, create a new migration to update it!

## schema migrations

```/migrations/``` - sql for the schema updates.
These are sql statements, try to keep each migration to most atomic form as possible, eg one table per migration.

```/reverts/``` - sql files to revert any mirgation - for any mirgration file there should be a revert one.

For testing you should be able to run a migration then revert script then migration script again correctly.

Any lookup or global data which need to be in all enviroments should be added here.



# Running the DB update scripts

Running the schemaupdate-userperms_service.sh script example in Windows:

```
set PGPASSWORD=change_me_to_password
set WSLENV=PGPASSWORD
set WSLUSER=your_wsl_user_name  #Note only needed if running sqlc

bash schemaupdate-userperms_service.sh -h dev-ca-central-1-postgres.cluster-c5ikek2awdr4.ca-central-1.rds.amazonaws.com -u postgres -d userperms_service_ca_dev -m 0013-add-columns-and-transform-data.sql
```

Running the script example in Linux/OSX:
```
export PGPASSWORD=change_me_to_password
bash schemaupdate-userperms_service.sh -h dev-ca-central-1-postgres.cluster-c5ikek2awdr4.ca-central-1.rds.amazonaws.com -u postgres -d userperms_service_ca_dev -m 0013-add-columns-and-transform-data.sql
```


```
export DBUSER=postgres
export DB_NAME=userperms_service_ca_dev
export DBHOST=dev-ca-central-1-postgres.cluster-c5ikek2awdr4.ca-central-1.rds.amazonaws.com
export PGPASSWORD=postgrespass
export SKIPCONFIRM=yes

./schemaupdate-userperms_service.sh -h $DBHOST -u $DBUSER -d $DB_NAME -m 0001-add-table-db_schema_number_and_db_changes_log.sql
./schemaupdate-userperms_service.sh -h $DBHOST -u $DBUSER -d $DB_NAME -m 0002-add-role-equine_service_user.sql
...
```



If odd error like below:

```
schemaupdate-userperms_service.sh: line 2: $'\r': command not found
schemaupdate-userperms_service.sh: line 4: $'\r': command not found

schemaupdate-userperms_service.sh: line 6: $'\r': command not found
schemaupdate-userperms_service.sh: line 7: syntax error near unexpected token `$'{\r''
schemaupdate-userperms_service.sh: line 7: `function show_help {
```

or

```
bash dataupdate-userperms_service.sh -h %DBHOST% -u postgres -d userperms_service_ca_dev -e dev -f ER-2619.txt

 -- value: 1
 -- value: 1ncedata/sc_reference-object-equine-2022-05-30.sql
Cannot find: dataupdates/
```

Check script has not been edited in windows and had CR LF added to the line endings needs to just be LF.
Maybe run ```git config --global core.autocrlf false``` to fix and reclone repo.


# Machine setup to run scripts and dev info

In order to run the scripts in windows you will need to install Windows Subsystems for Linux and then Ubuntu or another linux from the Microsoft Store. 

Once Ubuntu/linux installed from bash prompt run ```sudo apt install postgresql-client-common``` to install psql cli client.

For OSX see: https://dba.stackexchange.com/questions/3005/how-to-run-psql-on-mac-os-x 






# Update Ubuntu for windows to postgres 13

Needed for testing or pg_dump wont work.

```
sudo apt update && sudo apt -y full-upgrade
sudo apt update
sudo apt install postgresql-client-13
sudo apt install postgresql-client
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc|sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |sudo tee  /etc/apt/sources.list.d/pgdg.list
sudo apt update
sudo apt install  postgresql-client-13
```




# Notes for automated testing schema population:

TO create db in docker image etc:

```createdb -h $DBHOST -p $DBPORT -U $DBUSER $DBNAME ```

Need to have PGPASSWORD set as env var  
DBUSER is server admin user usually "postgres"
DBPORT normally 5432
DBNAME like "equine_service_dev"



Something like the below run from ```/internal/adapters/secondary/postgres/``` should populate the db to the latest migrations:

```
echo "Populating schema"
export SKIPCONFIRM=yes
for migrationfile in `ls ./migrations/`; do
    echo "Running $migrationfile"
    ./schemaupdate-userperms_service.sh -h $DBHOST -u $DBUSER -d $DBNAME -m $migrationfile
    if [[ $? -ne 0 ]]; then 
        echo "Migration error"
        exit 1
    fi
done
echo "Populating schema complete"
```


After migration in test env/go docker images - maybe useful for to set password to something known:

```
echo "ALTER USER equine_service_user WITH PASSWORD 'new-password';" | psql -v ON_ERROR_STOP=1 -1qAt -h $DBHOST -U $DBUSER -d $DBNAME`
```