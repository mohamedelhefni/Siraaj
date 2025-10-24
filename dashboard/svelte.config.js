import adapter from '@sveltejs/adapter-static';
import { fileURLToPath } from 'url';
import path from 'path';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			// Output directory for the built files
			pages: path.resolve(__dirname, '../ui/dashboard'),
			assets: path.resolve(__dirname, '../ui/dashboard'),
			fallback: 'index.html',
			precompress: false,
			strict: true
		}),
		paths: {
			base: '/dashboard'
		}
	}
};

export default config;
