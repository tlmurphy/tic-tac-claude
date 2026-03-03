import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    // Output to ../dist so Go can embed it at compile time.
    outDir: "../dist",
    emptyOutDir: true,
  },
  server: {
    // Proxy WebSocket requests to the Go server during development.
    proxy: {
      "/ws": {
        target: "ws://localhost:8080",
        ws: true,
      },
    },
  },
});
