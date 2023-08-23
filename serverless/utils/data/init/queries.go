package init

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
  invitee varchar(255) NOT NULL,
  inviter varchar(255) NOT NULL,
  date_invited TIMESTAMP NOT NULL,
  FOREIGN KEY(invitee) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY(inviter) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE CASCADE
);`
