CREATE DATABASE trackdocs WITH OWNER = postgres ENCODING = 'UTF8' CONNECTION LIMIT = -1;
CREATE ROLE trackdocs WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION CONNECTION LIMIT -1 PASSWORD 'trackdocs';
GRANT ALL PRIVILEGES ON DATABASE trackdocs TO trackdocs;