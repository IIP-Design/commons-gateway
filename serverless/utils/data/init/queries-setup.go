package init

/*
THIS FILE IS LOCKED

This file defines the initial setup for the application database.
To ensure repeatable and recoverable updates, we no longer modify
these initial queries. Rather, any modifications should be completed
by means of migration queries that build sequentially upon this
foundation (and each other).

This formulation allows us to more easily trace schema changes and
revert changes if necessary.
*/

const teamsQuery = `CREATE TABLE IF NOT EXISTS teams (
  id VARCHAR(255) PRIMARY KEY,
  team_name VARCHAR(255) NOT NULL,
  active BOOLEAN NOT NULL,
  date_created TIMESTAMP NOT NULL
);`

const adminsQuery = `CREATE TABLE IF NOT EXISTS admins (
  email VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  team VARCHAR(255) NOT NULL,
  active BOOLEAN NOT NULL,
  date_created TIMESTAMP NOT NULL,
  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

const guestsQuery = `CREATE TABLE IF NOT EXISTS guests (
  email VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  team VARCHAR(255) NOT NULL,
  pass_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  expiration TIMESTAMP NOT NULL,
  date_created TIMESTAMP NOT NULL,
  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

const invitesQuery = `CREATE TABLE IF NOT EXISTS invites (
  invitee VARCHAR(255) NOT NULL,
  inviter VARCHAR(255) NOT NULL,
  date_invited TIMESTAMP NOT NULL,
  FOREIGN KEY(invitee) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY(inviter) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE CASCADE
);`

const uploadsQuery = `CREATE TABLE IF NOT EXISTS uploads (
  s3_id VARCHAR(24) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  team_id VARCHAR(255) NOT NULL,
  file_type VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT
);`

const migrationsQuery = `CREATE TABLE IF NOT EXISTS migrations (
	id VARCHAR(20) PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	date_applied TIMESTAMP NOT NULL
);`
