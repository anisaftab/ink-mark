import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, type Plugin } from 'vite';
import fs from 'node:fs';
import path from 'node:path';

/**
 * devLogPlugin — tees Vite warnings and errors to logs/vite.log.
 * Only activates when INKMARK_DEV_LOGS=1. Zero browser impact.
 */
function devLogPlugin(): Plugin {
	const enabled = process.env.INKMARK_DEV_LOGS === '1';
	if (!enabled) return { name: 'inkmark-devlog-noop' };

	const logsDir = path.resolve(process.cwd(), 'logs');
	fs.mkdirSync(logsDir, { recursive: true });
	const logPath = path.join(logsDir, 'vite.log');

	function writeLine(level: string, msg: string): void {
		const ts = new Date().toISOString();
		const line = JSON.stringify({ ts, level, source: 'vite', msg }) + '\n';
		fs.appendFileSync(logPath, line, 'utf8');
	}

	return {
		name: 'inkmark-devlog',
		configResolved() {
			writeLine('info', 'Vite dev server started');
		},
		buildStart() {
			writeLine('info', 'build start');
		},
		// Hook into Vite's logger by wrapping it after config is set.
		configureServer(server) {
			const original = server.config.logger;

			const wrap = (level: 'warn' | 'error', fn: (...args: unknown[]) => void) =>
				(msg: string, ...rest: unknown[]) => {
					fn.call(original, msg, ...rest);
					writeLine(level, msg);
				};

			// @ts-expect-error — Vite's logger type doesn't expose a setter but the
			// object is writable; we patch in place without breaking the reference.
			server.config.logger.warn = wrap('warn', original.warn);
			// @ts-expect-error
			server.config.logger.error = wrap('error', original.error);
		}
	};
}

export default defineConfig({
	plugins: [
		devLogPlugin(),
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		})
	]
});
