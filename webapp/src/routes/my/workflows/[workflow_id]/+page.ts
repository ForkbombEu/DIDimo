export const load = async ({ params }) => {
	const { workflow_id } = params;

	return {
		workflow_id
	};
};
