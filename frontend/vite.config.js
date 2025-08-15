import { fileURLToPath, URL } from "node:url";
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd());

    return {
        plugins: [vue(), vueDevTools(), tailwindcss()],
        server: {
            port: 5000,
            proxy: {
                "/api": {
                    target: env.VITE_BACKEND_BASE_URL,
                    changeOrigin: true,
                    secure: false,
                },
                "/websocket": {
                    target: env.VITE_BACKEND_BASE_URL,
                    ws: true,
                    changeOrigin: true,
                },
            },
        },
        resolve: {
            alias: {
                "@": fileURLToPath(new URL("./src", import.meta.url)),
            },
        },
    };
});
