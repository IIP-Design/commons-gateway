#!/bin/sh

container=gateway_db
database=gateway_dev
user=gateway_dev

truncate_query='TRUNCATE admins, credentials, invites;'

docker exec -it $container psql -d $database -U $user -c "$truncate_query"