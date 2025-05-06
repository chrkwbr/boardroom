create database boardroom;
create user boardroom with encrypted password 'boardroom';
grant all privileges on database boardroom to boardroom;
grant all privileges on schema public to boardroom;

-- Function to grant privileges on all existing schemas
DO $$
DECLARE
r RECORD;
BEGIN
FOR r IN (SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema'))
    LOOP
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ' || quote_ident(r.schema_name) || ' TO boardroom';
END LOOP;
END $$;

-- Function to grant privileges on newly created schemas
CREATE OR REPLACE FUNCTION grant_privileges_on_new_schema()
RETURNS event_trigger AS $$
BEGIN
EXECUTE 'GRANT ALL PRIVILEGES ON SCHEMA ' || current_schema() || ' TO boardroom';
END;
$$ LANGUAGE plpgsql;

-- Event trigger to call the function on schema creation
CREATE EVENT TRIGGER grant_privileges_on_new_schema_trigger
ON ddl_command_end
WHEN TAG IN ('CREATE SCHEMA')
EXECUTE PROCEDURE grant_privileges_on_new_schema();