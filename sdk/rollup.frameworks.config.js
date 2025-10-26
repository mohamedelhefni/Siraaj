import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import esbuild from 'rollup-plugin-esbuild';
import dts from 'rollup-plugin-dts';

const frameworks = ['react', 'vue', 'svelte', 'preact', 'next', 'nuxt'];

const configs = frameworks.flatMap((framework) => [
    // ESM and UMD builds
    {
        input: `src/frameworks/${framework}/index.ts`,
        output: [
            {
                file: `dist/${framework}.esm.js`,
                format: 'esm',
                sourcemap: true,
            },
            {
                file: `dist/${framework}.umd.js`,
                format: 'umd',
                name: `Siraaj${framework.charAt(0).toUpperCase() + framework.slice(1)}`,
                sourcemap: true,
                globals: {
                    react: 'React',
                    vue: 'Vue',
                    svelte: 'Svelte',
                    preact: 'preact',
                    'next/navigation': 'NextNavigation',
                    'vue-router': 'VueRouter',
                },
            },
        ],
        external: ['react', 'vue', 'svelte', 'preact', 'next/navigation', 'vue-router'],
        plugins: [
            resolve(),
            commonjs(),
            esbuild({
                minify: false,
                target: 'es2020',
                jsx: framework === 'react' || framework === 'preact' || framework === 'next' ? 'transform' : undefined,
                jsxFactory: framework === 'preact' ? 'h' : undefined,
            }),
        ],
    },
    // TypeScript declarations
    {
        input: `src/frameworks/${framework}/index.ts`,
        output: {
            file: `dist/${framework}.d.ts`,
            format: 'esm',
        },
        external: ['react', 'vue', 'svelte', 'preact', 'next/navigation', 'vue-router'],
        plugins: [dts()],
    },
]);

export default configs;
