import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getUserById, updateUser, toggleUserStatus, deleteUser } from '$lib/api/user';

export const load: PageServerLoad = async (event) => {
	const { params, fetch: eventFetch, parent } = event;
	await parent(); // Memastikan otentikasi dari layout induk selesai

	try {
		const response = await getUserById(params.id, eventFetch);
		return {
			user: response.data,
		};
	} catch (error: any) {
		console.error('Gagal memuat pengguna:', error);
		// Jika pengguna tidak ditemukan atau ada error lain, arahkan kembali ke daftar pengguna
		throw redirect(303, '/user');
	}
};

export const actions: Actions = {
	// Aksi untuk memperbarui data pengguna
	update: async ({ request, params, fetch: eventFetch }) => {
		const formData = await request.formData();
		const userData = {
			username: formData.get('username') as string,
			email: formData.get('email') as string,
			role: formData.get('role') as 'admin' | 'manager' | 'cashier',
		};

		try {
			await updateUser(params.id, userData, eventFetch);
			return { success: true, message: 'Data pengguna berhasil diperbarui.' };
		} catch (error: any) {
			return fail(400, { error: error.message });
		}
	},

	// Aksi untuk mengubah status
	toggleStatus: async ({ params, fetch: eventFetch }) => {
		try {
			await toggleUserStatus(params.id, eventFetch);
			return { success: true, message: 'Status pengguna berhasil diubah.' };
		} catch (error: any) {
			return fail(400, { error: error.message });
		}
	},

	// Aksi untuk menghapus pengguna
	delete: async ({ params, fetch: eventFetch }) => {
		try {
			await deleteUser(params.id, eventFetch);
		} catch (error: any) {
			return fail(400, { error: error.message });
		}
		// Jika berhasil, arahkan ke halaman daftar pengguna
		throw redirect(303, '/user');
	},
};
