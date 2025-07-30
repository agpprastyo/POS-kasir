import { redirect } from '@sveltejs/kit';
import type {  LayoutLoad } from './$types';
import { userProfile } from '$lib/stores';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

// Fungsi load ini akan berjalan sebelum layout dan halaman di dalamnya dirender.
export const load: LayoutLoad = async ({ fetch: eventFetch }) => {
	try {
		// Panggil API untuk memeriksa sesi. eventFetch akan meneruskan cookie.
		const response = await eventFetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/me`, {
			method: 'GET',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
		});

		console.log("Checking auth session...");
		console.log("Response status:", response.status);

		if (!response.ok) {
			// Jika tidak terotentikasi, lemparkan redirect.
			// Ini akan menghentikan rendering lebih lanjut dan mengarahkan ke login.
			throw redirect(303, '/login');
		}

		const result = await response.json();
		const profile = result.data;

		// Simpan profil ke store agar bisa diakses di mana saja.
		userProfile.set(profile);

		// Kembalikan data profil agar bisa diakses di +layout.svelte dan halaman turunannya.
		return { profile };

	} catch (error) {
		// Jika terjadi redirect atau error lain, tangkap dan arahkan ke login.
		if (error instanceof Response && error.status === 303) {
			throw error; // Lanjutkan redirect
		}
		console.error("Layout auth check failed:", error);
		throw redirect(303, '/login');
	}
};
