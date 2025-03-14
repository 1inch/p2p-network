import { defineConfig } from 'vite';

export default defineConfig({
  root: '',
  publicDir: 'public',
  build: {
    rollupOptions: {
      input: 'index.html',
    },
    outDir: 'dist'
  }
});
