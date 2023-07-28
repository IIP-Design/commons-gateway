# Readme

Serverless functions to grant temporary upload ability to Aprimo.

## Functions

### Provision Access

```mermaid
flowchart TD
  A[Receive guest and admin email addresses]
  B[Check if guest already has access]
  A --> B
  C[NO: Generate password, salt, and hash]
  D[YES: Notify admin that access already granted]
  B --> C
  B --> D
  E["Save data
    1. Guest password, salt, and email
    2. admin-guest association"]
  F["Send notifications
    1. Email guest their password
    2. Notify admin password sent"]
  C --> E
  E --> F
```

### Grant Access

```mermaid
flowchart TD
  A[Receive guest email address]
  B[Check if provided email is in the DB]
  A --> B
  C[YES: Send hash and salt to login app]
  D[NO: Notify guest they do not have access]
  B --> C
  B --> D
```

### Upload File(s)

## Todo

- [ ] Web portal for admin to add guest user email
- [ ] Okta authentication for admin portal
- [ ] Save guest user - admin relationship to DB
- [ ] Check if user already has access
- [x] Generate guest user password and seed
- [ ] Send guest user email with password
- [x] Hash the password salt combo
- [x] Save the hash and the salt to the DB with the user email
- [ ] Upload portal for guest user
- [ ] User inputs email and password
- [ ] Retrieve the salt and the hash from the DB
- [ ] Hash the password input with the retrieved salt
- [ ] Compare the hashes and if matched grant access
- [ ] Upload files from input to S3
- [ ] Clean up DB removing user entry with one-time password
