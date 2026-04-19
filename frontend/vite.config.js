import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
var src = function (segment) { return path.resolve(process.cwd(), 'src', segment); };
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            app: src('app'),
            components: src('components'),
            constants: src('constants'),
            features: src('features'),
            hooks: src('hooks'),
            pages: src('pages'),
            services: src('services'),
            store: src('store'),
            types: src('types'),
            utils: src('utils'),
        },
    },
    server: {
        port: 3000,
    },
});
