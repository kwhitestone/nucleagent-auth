import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "node:path";

// Dev server: http://localhost:6678
// API requests under /api are proxied to the prism-fusion backend on :6670
// so the browser never hits cross-origin limits during local development.
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": resolve(process.cwd(), "src"),
    },
  },
  server: {
    port: 6678,
    proxy: {
      "/api": {
        target: "http://localhost:6670",
        changeOrigin: true,
      },
    },
  },
});
