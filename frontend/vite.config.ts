import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // allow external access
    allowedHosts: ['frontend'], // allow NGINX to access via container hostname
    port: 3000,
  },
});