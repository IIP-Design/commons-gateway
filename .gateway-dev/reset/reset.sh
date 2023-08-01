#!/bin/sh

container=gateway_db
database=gateway_dev
user=gateway_dev

docker exec -it $container psql -d $database -U $user -c 'TRUNCATE credentials, invites;'