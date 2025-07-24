import { getProducts, getCategories } from '$lib/api/product';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch: eventFetch, parent }) => {
	// Memastikan otentikasi dari layout induk selesai
	await parent();

	// Ambil parameter dari URL untuk filter dan paginasi
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 12; // Default limit untuk tampilan grid
	const search = url.searchParams.get('search') || '';
	const category_id = url.searchParams.get('category_id') || '';
	const view = url.searchParams.get('view') || 'grid'; // Parameter untuk mode tampilan

	const queryParams = { page, limit, search, category_id };

	try {
		// Panggil kedua API secara bersamaan untuk efisiensi
		const [productsResponse, categoriesResponse] = await Promise.all([
			getProducts(queryParams, eventFetch),
			getCategories(eventFetch)
		]);

		return {
			// **FIX:** Pastikan products selalu berupa array, bahkan jika API mengembalikan null
			products: productsResponse.data.products || [],
			pagination: productsResponse.data.pagination,
			categories: categoriesResponse.data,
			queryParams,
			view, // Kirim mode tampilan ke UI
			error: null
		};
	} catch (error: any) {
		console.error('Error loading products:', error);
		return {
			products: [],
			pagination: null,
			categories: [],
			error: error.message,
			queryParams,
			view
		};
	}
};
