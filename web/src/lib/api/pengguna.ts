import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type { UserQueryParams, UsersApiResponse, CreateUserRequest, Profile } from '$lib/types';

/**
 * Mengambil daftar pengguna dari API dengan filter dan paginasi.
 * Menggunakan fetch kustom jika disediakan (untuk panggilan dari server SvelteKit).
 * @param {UserQueryParams} params - Objek yang berisi parameter query.
 * @param {typeof fetch} [customFetch=fetch] - Instance fetch opsional.
 * @returns {Promise<{ data: UsersApiResponse }>} - Respons dari API.
 */
export async function getUsers(params: UserQueryParams, customFetch: typeof fetch = fetch): Promise<{ data: UsersApiResponse }> {
	const query = new URLSearchParams();

	if (params.page) query.append('page', params.page.toString());
	if (params.limit) query.append('limit', params.limit.toString());
	if (params.search) query.append('search', params.search);
	if (params.role) query.append('role', params.role);
	if (params.is_active !== undefined && params.is_active !== '') query.append('is_active', String(params.is_active));
	if (params.sortBy) query.append('sortBy', params.sortBy);
	if (params.sortOrder) query.append('sortOrder', params.sortOrder);

	const response = await customFetch(`${PUBLIC_API_BASE_URL}/api/v1/users?${query.toString()}`, {
		method: 'GET',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
	});

	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || 'Gagal mengambil data pengguna.');
	}

	return result;
}

/**
 * Membuat pengguna baru.
 * @param {CreateUserRequest} userData - Data pengguna baru.
 * @returns {Promise<{ data: Profile }>} - Respons dari API.
 */
export async function createUser(userData: CreateUserRequest): Promise<{ data: Profile }> {
	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/users`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
		body: JSON.stringify(userData),
	});

	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || 'Gagal membuat pengguna baru.');
	}

	return result;
}
