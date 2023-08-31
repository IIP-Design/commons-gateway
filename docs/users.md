| Capability                 | Super Admin | Admin | Guest Admin | Guest |
| -------------------------- | ----------- | ----- | ----------- | ----- |
| upload file                | ✅          | ✅    | ✅          | ✅    |
| invite guest               | ✅          | ✅    | ❌          | ❌    |
| propose invitation         | ❌          | ❌    | ✅          | ❌    |
| accept proposed invitation | ✅          | ✅    | ❌          | ❌    |
| add/edit super admin user  | ✅          | ❌    | ❌          | ❌    |
| add/edit admin user        | ✅          | ❌    | ❌          | ❌    |
| add/edit guest admin user  | ✅          | ✅    | ❌          | ❌    |
| add/edit team              | ✅          | ❌    | ❌          | ❌    |

| Page Access            | Super Admin | Admin | Guest Admin | Guest |
| ---------------------- | ----------- | ----- | ----------- | ----- |
| upload                 | ✅          | ✅    | ✅          | ✅    |
| new user (to add)      | ✅          | ✅    | ❌          | ❌    |
| user (to edit)         | ✅          | ✅    | ❌          | ❌    |
| propose invite         | ❌          | ❌    | ✅          | ❌    |
| pending invite         | ✅          | ✅    | ❌          | ❌    |
| admins (manage admins) | ✅          | ❌    | ❌          | ❌    |
| teams (manage teams)   | ✅          | ❌    | ❌          | ❌    |

Auth

super admins & admins - Okta
guest admin & guest - password

## Suggested Schema Updates

**General Changes:**

1. Given that we are allowing updates to teams and users it may be wise to set ad date modified value for each item. The initial value of date modified can match the date created.
1. Update the on the delete actions for referenced admins and guests from `CASCADE` to `RESTRICT`. This should assist with record keeping and auditing, although it may grow the DB.

**Role-Specific Changes:**

Add an administrator role enumeration and use it to differentiate between the two principle admin types.

```diff
+CREATE TYPE ADMIN_ROLE AS ENUM ('super admin', 'admin');

CREATE TABLE IF NOT EXISTS admins (
  email VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
+  role ADMIN_ROLE NOT NULL DEFAULT 'admin',
  team VARCHAR(255) NOT NULL,
  active BOOLEAN NOT NULL,
  date_created TIMESTAMP NOT NULL,
+  date_modified TIMESTAMP NOT NULL,
-  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE
+  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

Similarly, add a guest role enumeration and use it to differentiate between the two guest types.

```diff
+CREATE TYPE GUEST_ROLE AS ENUM ('guest admin', 'guest');

CREATE TABLE IF NOT EXISTS guests (
  email VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
+  role GUEST_ROLE NOT NULL DEFAULT 'guest',
  team VARCHAR(255) NOT NULL,
  pass_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  expiration TIMESTAMP NOT NULL,
  date_created TIMESTAMP NOT NULL,
+  date_modified TIMESTAMP NOT NULL,
-  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE
+  FOREIGN KEY(team) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

Add an optional proposer column to the invites to indicate invites that originate from guest admins. Additionally, such invites will have a `pending` property set to `TRUE` until the invite is approved by an admin. Once approved, the pending value will become `FALSE` and `inviter` will be set to the user who approved the invitation. This change necessitates us removing the not null constraint on the inviter as they will be unknown until the approval occurs.

```diff
CREATE TABLE IF NOT EXISTS invites (
  invitee VARCHAR(255) NOT NULL,
-  inviter VARCHAR(255) NOT NULL,
+  inviter VARCHAR(255),
+  proposer VARCHAR(255),
  date_invited TIMESTAMP NOT NULL,
+ pending BOOLEAN NOT NULL DEFAULT FALSE,
-  FOREIGN KEY(invitee) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE,
-  FOREIGN KEY(inviter) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE CASCADE
+  FOREIGN KEY(invitee) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE RESTRICT,
+  FOREIGN KEY(inviter) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE RESTRICT,
+  FOREIGN KEY(proposer) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

**Outstanding Questions:**

How can we set foreign key constraints on the user id in the uploads table when we have both admins and guests can upload data.

```sql
CREATE TABLE IF NOT EXISTS uploads (
  s3_id VARCHAR(24) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL, -- Would love for this to be a foreign key
  team_id VARCHAR(255) NOT NULL,
  file_type VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

```diff
+CREATE TABLE IF NOT EXISTS all_users (
+  user_id VARCHAR(24) PRIMARY KEY,
+  admin_id VARCHAR(255),
+  guest_id VARCHAR(255),
+  FOREIGN KEY(admin_id) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE CASCADE
+  FOREIGN KEY(guest_id) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE
+);

CREATE TABLE IF NOT EXISTS uploads (
  s3_id VARCHAR(24) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  team_id VARCHAR(255) NOT NULL,
  file_type VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT
+  FOREIGN KEY(user_id) REFERENCES all_users(user_id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```
