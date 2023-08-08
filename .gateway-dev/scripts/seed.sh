#!/bin/sh

# DB connection details.
container=gateway_db
database=gateway_dev
user=gateway_dev

# Set the data to populate the dev database with.
current_time=\'$(date +"%Y-%m-%d %T")\'
admin=\'test@gmail.com\'
guest=\'test1@gmail.com\'
hash=\'i82y9CY9olYWVDP3BPwdK1lVhBv60FEo3UtIbSJO8zQ=\'
password=kaW6PfO3MHZto2MlNQMV
salt=\'Aez71qglib\'

# Formulate the seeding queries.
admins_query="INSERT INTO admins (email, active, date_created) VALUES ($admin, true, $current_time);"
creds_query="INSERT INTO credentials (email, pass_hash, salt, date_created) VALUES ($guest, $hash, $salt, $current_time);"
invites_query="INSERT INTO invites (invitee, inviter, date_invited) VALUES ($admin, $guest, $current_time);"

# Run the queries.
echo "\nCreating admin user..."
docker exec -it $container psql -d $database -U $user -c "$admins_query"

echo "\nCreating guest user..."
docker exec -it $container psql -d $database -U $user -c "$creds_query"

echo "\nSimulating the guest user's invitation..."
docker exec -it $container psql -d $database -U $user -c "$invites_query"

echo "\nSEEDING COMPLETE! Added the admin user $admin and a guest user $guest who has the password $password."