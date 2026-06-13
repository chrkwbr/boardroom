import { defineConfig } from "vite";
import deno from "@deno/vite-plugin";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [deno(), react(), tailwindcss()],
  server: {
    proxy: {
      "/api/command": {
        target: "http://localhost:8088",
        changeOrigin: true,
      },
      "/api/query": {
        target: "http://localhost:8088",
        changeOrigin: true,
      },
      "/ws": {
        target: "ws://localhost:8088",
        ws: true,
      },
    },
  },
});
