#!/usr/bin/env node

import { minify } from 'terser';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join, extname } from 'path';

const DIST_DIR = 'dist';

const terserOptions = {
    compress: {
        passes: 3,
        pure_getters: true,
        unsafe: true,
        unsafe_comps: true,
        unsafe_math: true,
        unsafe_methods: true,
        drop_console: true,
        drop_debugger: true,
        ecma: 2020,
    },
    mangle: {
        properties: {
            regex: /^_/,
        },
    },
    format: {
        comments: false,
        ecma: 2020,
    },
    sourceMap: {
        filename: '',
        url: '',
    },
    ecma: 2020,
};

async function minifyFile(filePath) {
    const code = readFileSync(filePath, 'utf-8');
    const sourceMapPath = `${filePath}.map`;
    const hasSourceMap = statSync(sourceMapPath, { throwIfNoEntry: false });

    const outputPath = filePath.replace(/\.js$/, '.min.js');
    const sourceMapOutputPath = `${outputPath}.map`;

    const options = { ...terserOptions };

    if (hasSourceMap) {
        options.sourceMap = {
            content: readFileSync(sourceMapPath, 'utf-8'),
            filename: outputPath.split('/').pop(),
            url: `${outputPath.split('/').pop()}.map`,
        };
    }

    try {
        const result = await minify(code, options);

        if (!result.code) {
            console.error(`❌ Failed to minify ${filePath}: No output`);
            return;
        }

        writeFileSync(outputPath, result.code, 'utf-8');

        if (result.map && hasSourceMap) {
            writeFileSync(sourceMapOutputPath, result.map, 'utf-8');
        }

        const originalSize = (code.length / 1024).toFixed(2);
        const minifiedSize = (result.code.length / 1024).toFixed(2);
        const savings = (((code.length - result.code.length) / code.length) * 100).toFixed(1);

        console.log(`✅ ${filePath} -> ${outputPath}`);
        console.log(`   ${originalSize} KB -> ${minifiedSize} KB (${savings}% reduction)`);
    } catch (error) {
        console.error(`❌ Error minifying ${filePath}:`, error.message);
    }
}

async function processDirectory(dir) {
    const files = readdirSync(dir);

    for (const file of files) {
        const filePath = join(dir, file);
        const stat = statSync(filePath);

        if (stat.isDirectory()) {
            continue; // Skip directories
        }

        if (extname(file) === '.js' && !file.includes('.min.') && !file.includes('.map')) {
            await minifyFile(filePath);
        }
    }
}

console.log('🔨 Starting minification process...\n');

processDirectory(DIST_DIR)
    .then(() => {
        console.log('\n✨ Minification complete!');
    })
    .catch((error) => {
        console.error('❌ Minification failed:', error);
        process.exit(1);
    });
