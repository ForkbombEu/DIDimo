export function getExceptionMessage(e: unknown): string {
	if (e instanceof Error) {
		return e.message;
	} else {
		return JSON.stringify(e);
	}
}

//

export class NotBrowserError extends Error {}
