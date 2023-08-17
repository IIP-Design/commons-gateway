#!/bin/sh

# DB connection details.
container=gateway_db
database=gateway_dev
user=gateway_dev

# Query to delete all the dev data.
truncate_query='TRUNCATE admins, guests, invites, teams;'

# Run the queries.
echo "\nResetting the database..."
docker exec -it $container psql -d $database -U $user -c "$truncate_query"