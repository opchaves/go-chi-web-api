import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const envDir = path.resolve(process.cwd(), "..");
  const env = loadEnv(mode, envDir, "");

  return {
    base: env.VITE_BASE_URL,
    server: {
      port: +env.WEB_PORT,
      proxy: {
        "/api": {
          target: env.VITE_API_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ""),
        },
      },
    },
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
        '~bootstrap': path.resolve(__dirname, 'node_modules/bootstrap'),
      },
    },
  };
});
