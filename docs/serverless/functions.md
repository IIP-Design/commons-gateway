---
layout: page
title: Serverless Functions
---

This page serves as a reference for all of the Lambda functions that handle the data flow through out the application.

## Retrieve Admin Users

TODO

This operation returns a list of all users authorized to invite guest users.

## Retrieve Teams

TODO

This operation returns a list of all teams.

## Retrieve Guest Users

TODO

This operation returns a list of all users authorized for guest uploading.

## Admin User Create

This operation authorizes the provided email to add new guest users.

<div class="mermaid">
flowchart TD
  A[Receive new admin email address]
  B[Check if user is already an admin]
  A --> B
  C[NO: Add user to the list of admins]
  D[YES: Inform the user they are already an admin]
  B --> C
  B --> D
</div>

## Admin User Edit

TODO

This operation updates information for an existing admin user.

## Team Create

TODO

This operation creates a new Commons team.

## Provision Credentials

This operations generates and securely stores the temporary password for the

<div class="mermaid">
flowchart TD
  A[Receive guest and admin email addresses]
  B[Check if guest already has credentials]
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
</div>

## Retrieve Salt

This operation retrieves the salt used when generating a user's password salt.

<div class="mermaid">
flowchart TD
  A[Receive guest email address]
  B[Check if provided email is in the DB]
  A --> B
  C[YES: Send salt to login app]
  D[NO: Notify guest they do not have access]
  B --> C
  B --> D
</div>

## Upload File(s)

TODO

This operation handles the uploading of user files to an S3 bucket where they are scanned for malware.

## Send Notification

TODO

This operation sends an email notification.

<script src="https://cdnjs.cloudflare.com/ajax/libs/mermaid/10.3.1/mermaid.min.js"></script>

<script>
var config = {
  startOnLoad:true,
  theme: 'neutral',
  flowchart:{
    useMaxWidth:false,
    htmlLabels:true
  }
};

mermaid.initialize(config);
window.mermaid.init(undefined, document.querySelectorAll('.language-mermaid'));
</script>
