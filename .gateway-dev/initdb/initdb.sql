CREATE TABLE IF NOT EXISTS otp (
  email VARCHAR(255) UNIQUE NOT NULL,
  otp_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  date_created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS invites (
  invitee varchar(255) NOT NULL,
  inviter varchar(255) NOT NULL,
  date_invited TIMESTAMP NOT NULL
);