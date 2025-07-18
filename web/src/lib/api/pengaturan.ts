import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type { Profile, UpdatePasswordRequest } from '$lib/types';

/**
 * Fungsi pembantu untuk menangani respons dari API.
 * @param response - Objek Response dari fetch.
 * @returns Promise yang resolve dengan data JSON.
 * @throws Error jika respons tidak ok.
 */
async function handleResponse(response: Response) {
	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || `HTTP error! status: ${response.status}`);
	}
	return result;
}

/**
 * Mengambil data profil pengguna yang sedang login.
 * @returns {Promise<{ data: Profile }>} - Objek yang berisi profil pengguna.
 */
export async function getMyProfile(): Promise<{ data: Profile }> {
	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/me`, {
		method: 'GET',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
	});
	return handleResponse(response);
}

/**
 * Mengunggah dan memperbarui avatar pengguna.
 * @param {File} avatarFile - File gambar yang akan diunggah.
 * @returns {Promise<{ data: Profile }>} - Objek yang berisi profil pengguna yang sudah diperbarui.
 */
export async function updateAvatar(avatarFile: File): Promise<{ data: Profile }> {
	const formData = new FormData();
	formData.append('avatar', avatarFile);

	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/me/avatar`, {
		method: 'PUT',
		credentials: 'include',
		body: formData, // Tidak perlu set Content-Type, browser akan menanganinya untuk FormData
	});
	return handleResponse(response);
}

/**
 * Memperbarui password pengguna.
 * @param {UpdatePasswordRequest} passwordData - Data password lama dan baru.
 * @returns {Promise<any>} - Respons sukses dari API.
 */
export async function updatePassword(passwordData: UpdatePasswordRequest): Promise<any> {
	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/me/password`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
		body: JSON.stringify(passwordData),
	});
	return handleResponse(response);
}
