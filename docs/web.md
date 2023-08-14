---
layout: page
title: Client Application
---

The client application is how the user interfaces with the Aprimo API. The primary purpose of the site is to provide guest users the ability to upload files to Content Commons without provisioning them full access. The interface also as an admin portal where authorized team members can invite guest users, as well as a super admin page where authorized users can create/edit team and add new admin users.

This interface is deployed as a statically generated site built on [Astro](https://astro.build/). Data interactions are processed offsite using a series of Lambda [serverless functions]({{ '/functions' | relative_url }}). These Lambda functions handle authentication, fetching and updating system data, and the connection to Aprimo.

## Project Structure

The client application is found inside of the web directory and is a fairly standard Astro project. You'll find the following folders and files:

```bash
web/
│
├── public/ # Unprocessed static assets (images, fonts, icons, etc.)
│   └── dos_seal.svg # Used as a favicon
│
├── src/
│   ├── components/ # Reusable UI components
│   │   └── Button.astro # Astro components
│   │   └── Table.tsx # React components
│   │
│   ├── layouts/ # Shared UI configuration used by multiple pages
│   ├── pages/ # Website page files
│   ├── styles/ # Application-wide style sheet
│   └── utils/ # Reusable utility functions
│
└── package.json
```

Astro looks for `.astro` or `.md` files in the `src/pages/` directory. Each page is exposed as a route based on its file name. The page `404.astro` serves as a fallback should the user navigate to an non-existent page.

Most of the site is built of simple Astro components, however, where we need to maintain client-side state we can use React components. These components can be import just as any other Astro components and utilized as a standard React component would be:

```js
// src/pages/somepage.astro
---
import MyReactComponent from '../components/MyReactComponent.tsx';
---

<div>
  <h1>Use React components directly in Astro!</h1>
  <MyReactComponent value="my-prop"/>
</div>
```

## Commands

All commands are run from the root of the project (i.e., the `web` directory), from a terminal:

| Command                   | Action                                           |
| :------------------------ | :----------------------------------------------- |
| `npm run dev`             | Starts local dev server at `localhost:3000`      |
| `npm run build`           | Build your production site to `./dist/`          |
| `npm run preview`         | Preview your build locally, before deploying     |
| `npm run astro ...`       | Run CLI commands like `astro add`, `astro check` |
| `npm run astro -- --help` | Get help using the Astro CLI                     |

## More on Astro

For more information on Astro, please refer to their [official documentation](https://docs.astro.build).

## Mockups

Below you will find mockups for the user flow for each of the user types.

- [Super Admin User](https://preview.uxpin.com/e0f263260e840c070f1c9796be483a6741acb510#/pages//simulate/sitemap?mode=i)
- [Admin User](https://preview.uxpin.com/26933858eb418523b89317817ab80570867d7dc3#/pages//simulate/sitemap?mode=i)
- [Guest User Team Lead](https://preview.uxpin.com/3574e3fc7f3e3d9658f5757a67d7f9e7e53c8f2b#/pages//simulate/sitemap?mode=i)
- [Guest User](https://preview.uxpin.com/c690dff0b5961a53b47e1ba62141632457175178#/pages//simulate/sitemap?mode=i)
