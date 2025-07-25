import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type { Category, CategoryRequest, CategoryWithCount } from '$lib/types';

/**
 * Fungsi pembantu untuk menangani respons dari API.
 */
async function handleResponse(response: Response) {
	if (response.status === 204 || response.status === 200 && response.headers.get('content-length') === '0') {
		return; // Handle No Content response for delete/update
	}
	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || `HTTP error! status: ${response.status}`);
	}
	return result;
}

/**
 * Mengambil daftar kategori beserta jumlah produk.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: CategoryWithCount[] }>} - Daftar kategori.
 */
export async function getCategoriesWithCount(customFetch: typeof fetch): Promise<{ data: CategoryWithCount[] }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/categories/count`);
	return handleResponse(response);
}

/**
 * Membuat kategori baru.
 * @param {CategoryRequest} categoryData - Nama kategori baru.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: Category }>} - Kategori yang baru dibuat.
 */
export async function createCategory(categoryData: CategoryRequest, customFetch: typeof fetch): Promise<{ data: Category }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/categories`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(categoryData),
	});
	return handleResponse(response);
}

/**
 * Memperbarui kategori.
 * @param {number} id - ID kategori.
 * @param {CategoryRequest} categoryData - Nama kategori yang baru.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 */
export async function updateCategory(id: number, categoryData: CategoryRequest, customFetch: typeof fetch): Promise<any> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/categories/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(categoryData),
	});
	return handleResponse(response);
}

/**
 * Menghapus kategori.
 * @param {number} id - ID kategori.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 */
export async function deleteCategory(id: number, customFetch: typeof fetch): Promise<any> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/categories/${id}`, {
		method: 'DELETE',
	});
	return handleResponse(response);
}
