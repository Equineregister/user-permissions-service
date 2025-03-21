DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'userperms_service_user') THEN
      RAISE NOTICE 'Role "userperms_service_user" already exists. Skipping.';
   ELSE
      CREATE ROLE userperms_service_user LOGIN PASSWORD 'change-ME';
   END IF;
END
$do$;
