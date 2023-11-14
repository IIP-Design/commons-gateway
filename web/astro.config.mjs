import { defineConfig } from 'astro/config'; // eslint-disable-line node/no-unpublished-import
import react from '@astrojs/react';

// https://astro.build/config
export default defineConfig( {
  integrations: [react()],
  scopedStyleStrategy: 'where',
  vite: {
    optimizeDeps: {
      exclude: ['date-fns'],
    },
  },
} );
