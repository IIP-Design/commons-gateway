#!/bin/sh

# DB connection details.
container=gateway_db
database=gateway_dev
user=gateway_dev

# Set the data to populate the dev database with.
current_time=$(date +"%Y-%m-%d %T")
expiration="2024-02-29 16:21:42"
admin=("alice@testmail.com" "Alice" "Apple" "1" true)
guest=("bob@testmail.com" "Bob" "Banana" "1")

# hash, password, salt
creds=("i82y9CY9olYWVDP3BPwdK1lVhBv60FEo3UtIbSJO8zQ=" "kaW6PfO3MHZto2MlNQMV" "Aez71qglib")

# Formulate the seeding queries.
teams_query="INSERT INTO teams (id, team_name, active,  date_created) VALUES ('1', 'Team Number One', true, '$current_time');"
admins_query="INSERT INTO admins (email, first_name, last_name, team, active, date_created) VALUES ('${admin[0]}', '${admin[1]}', '${admin[2]}', '${admin[3]}', '${admin[4]}', '$current_time');"
guest_query="INSERT INTO guests (email, first_name, last_name, team, pass_hash, salt, expiration, date_created) VALUES ('${guest[0]}', '${guest[1]}', '${guest[2]}', '${guest[3]}', '${creds[0]}', '${creds[2]}', '$expiration', '$current_time');"
invites_query="INSERT INTO invites (invitee, inviter, date_invited) VALUES ('${guest[0]}', '${admin[0]}', '$current_time');"

# Run the queries.
echo "\nCreating team..."
docker exec -it $container psql -d $database -U $user -c "$teams_query"

echo "\nCreating admin user..."
docker exec -it $container psql -d $database -U $user -c "$admins_query"

echo "\nCreating guest user..."
docker exec -it $container psql -d $database -U $user -c "$guest_query"

echo "\nSimulating the guest user's invitation..."
docker exec -it $container psql -d $database -U $user -c "$invites_query"

echo "\nSEEDING COMPLETE! Added the admin user $admin and a guest user $guest who has the password ${creds[1]}."