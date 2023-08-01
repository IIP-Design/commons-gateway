CREATE TABLE IF NOT EXISTS admins (
  email VARCHAR(255) UNIQUE NOT NULL,
  active BOOLEAN NOT NULL,
  date_created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS credentials (
  email VARCHAR(255) UNIQUE NOT NULL,
  pass_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  date_created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS invites (
  invitee varchar(255) NOT NULL,
  inviter varchar(255) NOT NULL,
  date_invited TIMESTAMP NOT NULL
);