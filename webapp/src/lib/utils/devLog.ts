export function log(...args: unknown[]) {
	if (process.env.NODE_ENV === 'development') {
		console.log(...args);
	}
}
