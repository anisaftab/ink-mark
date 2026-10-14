/**
 * devlog.ts — Agent-loop console patcher.
 *
 * When INKMARK_DEV_LOGS=1 is set on the Go side, this module intercepts
 * console.error, console.warn, and unhandled exceptions then forwards them
 * to the Go backend via LogBrowserEvent, which appends them to logs/console.log.
 *
 * Activation: imported once in +layout.svelte. Self-activates only when the
 * Wails bridge is present (i.e. running inside the app, not raw Vite/browser).
 * In raw Vite the window.go binding doesn't exist, so every call is a no-op.
 */

// Wails injects window.go at runtime — not available during SSR/Vite-only.
declare global {
	interface Window {
		go?: {
			main?: {
				App?: {
					LogBrowserEvent?: (level: string, message: string, stack: string) => Promise<void>;
				};
			};
		};
	}
}

function send(level: string, message: string, stack: string): void {
	window.go?.main?.App?.LogBrowserEvent?.(level, message, stack);
}

function stringify(args: unknown[]): string {
	return args
		.map((a) => {
			if (a instanceof Error) return a.message;
			if (typeof a === 'object') {
				try {
					return JSON.stringify(a);
				} catch {
					return String(a);
				}
			}
			return String(a);
		})
		.join(' ');
}

// Patch console.error
const _error = console.error.bind(console);
console.error = (...args: unknown[]) => {
	_error(...args);
	const first = args[0];
	const stack = first instanceof Error ? (first.stack ?? '') : '';
	send('error', stringify(args), stack);
};

// Patch console.warn
const _warn = console.warn.bind(console);
console.warn = (...args: unknown[]) => {
	_warn(...args);
	send('warn', stringify(args), '');
};

// Catch unhandled promise rejections
window.addEventListener('unhandledrejection', (e: PromiseRejectionEvent) => {
	const reason = e.reason;
	const message = reason instanceof Error ? reason.message : String(reason);
	const stack = reason instanceof Error ? (reason.stack ?? '') : '';
	send('error', `Unhandled rejection: ${message}`, stack);
});

// Catch uncaught errors
window.addEventListener('error', (e: ErrorEvent) => {
	send('error', e.message, e.error?.stack ?? '');
});
