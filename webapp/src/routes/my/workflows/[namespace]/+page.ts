export const load = async ({ params }) => {
	const { namespace } = params;

	return {
		namespace
	};
};
