import { getProducts, getCategories } from '$lib/api/product';
import type { PageLoad } from '../../../../.svelte-kit/types/src/routes';

export const load: PageLoad = async ({ url, fetch: eventFetch, parent }) => {
	await parent(); // Memastikan otentikasi

	const search = url.searchParams.get('search') || '';
	const category_id = url.searchParams.get('category_id') || '';

	// Muat semua produk (atau dengan paginasi yang besar) untuk ditampilkan
	const queryParams = { limit: 100, search, category_id };

	try {
		const [productsResponse, categoriesResponse] = await Promise.all([
			getProducts(queryParams, eventFetch),
			getCategories(eventFetch)
		]);

		return {
			products: productsResponse.data.products || [],
			categories: categoriesResponse.data,
			queryParams,
			error: null
		};
	} catch (error: any) {
		console.error('Error loading cashier data:', error);
		return {
			products: [],
			categories: [],
			error: error.message,
			queryParams
		};
	}
};
