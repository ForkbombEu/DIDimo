import z from 'zod';

//

export function zodFileSchema(options: { mimeTypes?: string[]; maxSize?: number } = {}) {
	const { mimeTypes, maxSize } = options;
	let s = z.instanceof(File);
	if (mimeTypes && mimeTypes.length > 0) {
		const mimes = mimeTypes as readonly string[];
		s = s.refine((file) => mimes.includes(file.type), `File type not: ${mimes.join(', ')}`);
	}
	if (maxSize) {
		s.refine((file) => file.size < maxSize, `File size bigger than ${maxSize} bytes`);
	}
	return s;
}

//

export function readFileAsBase64(file: File): Promise<string> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.readAsDataURL(file);
		reader.onload = () => {
			const base64string = (reader.result as string).split(',')[1];
			resolve(base64string);
		};
		reader.onerror = () => {
			reject(reader.error);
		};
	});
}

export function readFileAsString(file: File): Promise<string> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.readAsText(file);
		reader.onload = () => {
			const result = (reader.result as string).trim();
			resolve(result);
		};
		reader.onerror = () => {
			reject(reader.error);
		};
	});
}
