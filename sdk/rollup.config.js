import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import esbuild from 'rollup-plugin-esbuild';
import dts from 'rollup-plugin-dts';

const production = !process.env.ROLLUP_WATCH;

export default [
    // Core library (vanilla JS)
    {
        input: 'src/core/index.ts',
        output: [
            {
                file: 'dist/analytics.esm.js',
                format: 'esm',
                sourcemap: true,
            },
            {
                file: 'dist/analytics.umd.js',
                format: 'umd',
                name: 'SiraajAnalytics',
                sourcemap: true,
            },
        ],
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
