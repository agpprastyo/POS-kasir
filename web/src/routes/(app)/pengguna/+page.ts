import { getUsers } from '$lib/api/pengguna';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch: eventFetch }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const search = url.searchParams.get('search') || '';
	const role = (url.searchParams.get('role') as any) || '';
	const isActiveParam = url.searchParams.get('is_active');
	const is_active = isActiveParam === null ? '' : isActiveParam === 'true';

	const queryParams = { page, limit, search, role, is_active };

	try {
		// Panggil API dengan meneruskan `eventFetch`
		const response = await getUsers(queryParams, eventFetch);
		return {
			users: response.data.users,
			pagination: response.data.pagination,
			queryParams: queryParams
		};
	} catch (error: any) {
		console.error('Error loading users:', error);
		return {
			users: [],
			pagination: null,
			error: error.message,
			queryParams: queryParams
		};
	}
};
