---
layout: page
title: Commons Gateway
---

Commons Gateway portal that allows authorized guest users to upload files to Content Commons using the Aprimo API. Guest users are invited by site admins and provisioned with a one-time password as demonstrated in the diagram below.

![A chart displaying demonstrating a simplified data flow and infrastructural setup for the Commons Gateway application.]({{ '/assets/architectural-diagram.png' | relative_url }})

The project is composed to two parts:

1. **[Client Application:]({{ '/web' | relative_url }})** The `web` directory where the client application is found. This simple static site is built using [Astro](https://astro.build/) and serves as the interface between the user and the Aprimo API.
1. **[Serverless Functions:]({{ '/functions' | relative_url }})** The `serverless` directory, which contains a number of serverless functions that manage the data in the application.
