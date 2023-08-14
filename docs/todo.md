---
layout: page
title: Todo List
---

- [ ] **Web App - Page for super-admins to manage the application**
  - [ ] Okta authentication for super-admin portal
  - [ ] Retrieve the list of existing admin users (Retrieve Admin Users Lambda)
  - [ ] Retrieve the list of existing teams (Retrieve Teams Lambda)
- [ ] **Lambda Function - Retrieve Admin Users**
  - [ ] Return all registered admin users
- [ ] **Lambda Function - Retrieve Teams**
  - [ ] Return all registered teams
- [ ] **Web App - Page for admin to add guest user email**
  - [ ] Okta authentication for admin portal
  - [ ] Retrieve the list of existing guest users (Retrieve Guest Users Lambda)
  - [ ] Populate a table of users with the retrieved data
  - [ ] Allow for searching users
  - [ ] Clicking on the `New User` button opens a user form modal
  - [ ] Clicking on an existing user allows the admin to edit an existing user
  - [ ] Validate form inputs
  - [ ] Submitting the user form saves the user data and provides guest user access (Provision Credentials Lambda)
  - Allow admin to prematurely deactivate an existing guest user (Deactivation Guest User Lambda)
- [ ] **Lambda Function - Retrieve Guest Users**
  - [ ] Receive input from Web App admin page - team name
  - [ ] Return all users associated with the provided team
- [ ] **Lambda Function - Provision Credentials**
  - [x] Receive input from Web App admin page - user name
  - [x] Check if provided user email already has credentials
  - [ ] If so notify admin (else proceed as below)
  - [x] Save guest user - admin relationship to DB
  - [x] Generate guest user password and salt
  - [x] Hash the password salt combo
  - [x] Save the hash and the salt to the DB with the user email
  - [ ] Send guest user email with password
  - [ ] Notify admin inviter that password has been sent
- [ ] **Web App - Guest upload portal**
  - [x] User inputs email and password
  - [x] Retrieve the salt from the DB (Grant Access Lambda)
  - [x] Locally, hash the password input with the retrieved salt
  - [ ] Check for access cookie
  - [ ] Allow authenticated user to upload documents (Upload Files Lambda)
  - [ ] Record access?
  - [ ] On completed upload, remove user access (Cleanup Lambda)
- [ ] **Lambda Function - Grant Access**
  - [x] Receive input from Web App guest page
  - [x] Check if provided user email has access
  - [ ] If not, notify guest they do not have access (else proceed as below)
  - [x] Send the guest salt to the Web App guest page
  - [x] Receive the client generated hash
  - [x] Compare the received and stored hashes
  - [x] If the hashes do not match forbid access
  - [ ] If the hashes do match, return a session cookie granting access
- [ ] **Lambda Function - Upload files**
  - [ ] Upload files from input to S3
- [ ] **Lambda Function - Cleanup**
  - [ ] Receive input from Web App guest page
  - [ ] Removing user from access list
- [ ] **Lambda Function - Publish Upload to Aprimo**
- [ ] **Lambda Function - Add New Admin User**
  - [ ] Receive input from Web App admin page
  - [x] Check if provided user email belongs to an existing user
  - [x] Save user to the list of admins
- [ ] **Lambda Function - Deactivation Guest User**
  - [ ] Receive input from Web App admin page
  - [ ] Check if provided user email is indeed an admin
  - [ ] Set the `active` status on their entry to `false`

<style>
  .task-list {
    list-style: none;
    padding-inline-start: 1.5rem;
  }

  .task-list-item-checkbox {
    margin-right: 0.5rem;
  }
</style>
