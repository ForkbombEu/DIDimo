import z from 'zod';

const urlRegex =
	/^(https?:\/\/(?:www\.)?)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$/;

export const schema = z.object({ url: z.string().regex(urlRegex, 'Invalid URL') });
