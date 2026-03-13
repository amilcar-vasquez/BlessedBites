import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const API_TARGET = process.env.API_TARGET || 'http://localhost:8080';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    proxy: {
      // Forward /api/* to the Go API in local dev
      '/api': { target: API_TARGET, changeOrigin: true },
      // Forward /uploads/* to the Go API (which proxies to the static dir)
      // so dev and production behave identically
      '/uploads': { target: API_TARGET, changeOrigin: true }
    }
  }
});
