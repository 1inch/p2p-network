import { defineConfig } from 'vite';
import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import { createCandidateRouter } from '../../router';

export default defineConfig({
  root: '',
  publicDir: 'public',
  build: {
    rollupOptions: {
      input: 'index.html',
    },
    outDir: 'dist'
  },
  plugins: [
    {
      name: 'candidate-plugin',
      configureServer(server) {
        const app = express();

        app.use(cors());
        app.use(bodyParser.json());

        app.use('/api', createCandidateRouter());

        server.middlewares.use(app);
      },
    },
  ],
});
