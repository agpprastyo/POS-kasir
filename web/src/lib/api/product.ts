import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type {
	ProductQueryParams,
	ProductsApiResponse,
	Category,
	CreateProductRequest,
	Product,
	ProductOption, UpdateProductOptionRequest, UpdateProductRequest
} from '$lib/types';

/**
 * Fungsi pembantu untuk menangani respons dari API.
 */
async function handleResponse(response: Response) {
	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || `HTTP error! status: ${response.status}`);
	}
	return result;
}

/**
 * Mengambil daftar produk dari API dengan filter dan paginasi.
 * @param {ProductQueryParams} params - Objek yang berisi parameter query.
 * @param {typeof fetch} [customFetch=fetch] - Instance fetch opsional untuk SSR.
 * @returns {Promise<{ data: ProductsApiResponse }>} - Respons dari API.
 */
export async function getProducts(params: ProductQueryParams, customFetch: typeof fetch = fetch): Promise<{ data: ProductsApiResponse }> {
	const query = new URLSearchParams();

	if (params.page) query.append('page', params.page.toString());
	if (params.limit) query.append('limit', params.limit.toString());
	if (params.search) query.append('search', params.search);
	if (params.category_id) query.append('category_id', String(params.category_id));

	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products?${query.toString()}`, {
		method: 'GET',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
	});

	return handleResponse(response);
}

/**
 * Mengambil semua kategori produk.
 * @param {typeof fetch} [customFetch=fetch] - Instance fetch opsional untuk SSR.
 * @returns {Promise<{ data: Category[] }>} - Daftar semua kategori.
 */
export async function getCategories(customFetch: typeof fetch = fetch): Promise<{ data: Category[] }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/categories`, {
		method: 'GET',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
	});

	return handleResponse(response);
}

/**
 * Membuat produk baru.
 * @param {CreateProductRequest} productData - Data produk baru.
 * @returns {Promise<{ data: Product }>} - Data produk yang baru dibuat.
 */
export async function createProduct(productData: CreateProductRequest): Promise<{ data: Product }> {
	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/products`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
		body: JSON.stringify(productData),
	});
	return handleResponse(response);
}

/**
 * Mengunggah gambar untuk produk utama.
 * @param {string} productId - ID produk.
 * @param {File} imageFile - File gambar.
 * @param {typeof fetch} [customFetch=fetch] - Instance fetch opsional.
 * @returns {Promise<{ data: Product }>} - Data produk yang sudah diperbarui.
 */
export async function uploadProductImage(productId: string, imageFile: File, customFetch: typeof fetch = fetch): Promise<{ data: Product }> {
	const formData = new FormData();
	formData.append('image', imageFile);

	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${productId}/image`, {
		method: 'POST',
		body: formData,
	});
	return handleResponse(response);
}

/**
 * Mengunggah gambar untuk opsi produk.
 * @param {string} productId - ID produk.
 * @param {string} optionId - ID opsi produk.
 * @param {File} imageFile - File gambar.
 * @param {typeof fetch} [customFetch=fetch] - Instance fetch opsional.
 * @returns {Promise<{ data: ProductOption }>} - Data opsi produk yang sudah diperbarui.
 */
export async function uploadProductOptionImage(productId: string, optionId: string, imageFile: File, customFetch: typeof fetch = fetch): Promise<{ data: ProductOption }> {
	const formData = new FormData();
	formData.append('image', imageFile);

	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${productId}/options/${optionId}/image`, {
		method: 'POST',
		body: formData,
	});
	return handleResponse(response);
}

/**
 * Mengambil data satu produk berdasarkan ID.
 * @param {string} id - ID produk.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: Product }>} - Data detail produk.
 */
export async function getProductById(id: string, customFetch: typeof fetch): Promise<{ data: Product }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${id}`);
	return handleResponse(response);
}

/**
 * Memperbarui data produk utama.
 * @param {string} id - ID produk.
 * @param {UpdateProductRequest} productData - Data baru untuk produk.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: Product }>} - Data produk yang sudah diperbarui.
 */
export async function updateProduct(id: string, productData: UpdateProductRequest, customFetch: typeof fetch): Promise<{ data: Product }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${id}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(productData),
	});
	return handleResponse(response);
}

/**
 * Memperbarui data opsi produk.
 * @param {string} productId - ID produk.
 * @param {string} optionId - ID opsi produk.
 * @param {UpdateProductOptionRequest} optionData - Data baru untuk opsi.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: ProductOption }>} - Data opsi produk yang sudah diperbarui.
 */
export async function updateProductOption(productId: string, optionId: string, optionData: UpdateProductOptionRequest, customFetch: typeof fetch): Promise<{ data: ProductOption }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${productId}/options/${optionId}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(optionData),
	});
	return handleResponse(response);
}






/**
 * Membuat opsi produk baru.
 * @param {string} productId - ID produk.
 * @param {UpdateProductOptionRequest} optionData - Data opsi baru.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 * @returns {Promise<{ data: ProductOption }>} - Data opsi yang baru dibuat.
 */
export async function createProductOption(productId: string, optionData: UpdateProductOptionRequest, customFetch: typeof fetch): Promise<{ data: ProductOption }> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${productId}/options`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(optionData),
	});
	return handleResponse(response);
}

/**
 * Menghapus opsi produk.
 * @param {string} productId - ID produk.
 * @param {string} optionId - ID opsi yang akan dihapus.
 * @param {typeof fetch} customFetch - Instance fetch dari SvelteKit.
 */
export async function deleteProductOption(productId: string, optionId: string, customFetch: typeof fetch): Promise<any> {
	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/products/${productId}/options/${optionId}`, {
		method: 'DELETE',
	});
	return handleResponse(response);
}
