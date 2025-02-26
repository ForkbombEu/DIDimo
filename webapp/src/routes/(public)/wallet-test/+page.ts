import { error } from '@sveltejs/kit';

export const load = ({ url }) => {
	const qrContent = url.searchParams.get('qr');
	const workflowId = url.searchParams.get('workflow-id');
	if (!qrContent || !workflowId) error(404);

	return {
		qrContent,
		workflowId
	};
};
