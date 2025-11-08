import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import esbuild from 'rollup-plugin-esbuild';
import dts from 'rollup-plugin-dts';
import terser from '@rollup/plugin-terser';

const production = !process.env.ROLLUP_WATCH;

export default [
    // Core library (vanilla JS) - UMD for browser
    {
        input: 'src/core/index.ts',
        output: {
            file: 'analytics.js',
            format: 'umd',
            name: 'SiraajAnalytics',
            sourcemap: false,
            exports: 'named',
        },
        plugins: [
            resolve(),
            commonjs(),
            esbuild({
                minify: false,
                target: 'es2020',
            }),
        ],
    },
    // Minified version
    {
        input: 'src/core/index.ts',
        output: {
            file: 'analytics.min.js',
            format: 'umd',
            name: 'SiraajAnalytics',
            sourcemap: false,
            exports: 'named',
        },
        plugins: [
            resolve(),
            commonjs(),
            esbuild({
                minify: true,
                target: 'es2020',
            }),
            terser({
                compress: {
                    drop_console: true,
                    drop_debugger: true,
                },
            }),
        ],
    },
    // ES Module for modern bundlers
    {
        input: 'src/core/index.ts',
        output: {
            file: 'dist/analytics.esm.js',
            format: 'esm',
            sourcemap: true,
        },
        plugins: [
            resolve(),
            commonjs(),
            esbuild({
                minify: false,
                target: 'es2020',
            }),
        ],
    },
    // TypeScript declarations
    {
        input: 'src/core/index.ts',
        output: {
            file: 'dist/analytics.d.ts',
            format: 'esm',
        },
        plugins: [dts()],
    },
];
