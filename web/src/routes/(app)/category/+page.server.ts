import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getCategoriesWithCount, createCategory, updateCategory, deleteCategory } from '$lib/api/category';

export const load: PageServerLoad = async (event) => {
	await event.parent(); // Memastikan otentikasi

	try {
		const response = await getCategoriesWithCount(event.fetch);
		return {
			categories: response.data,
		};
	} catch (error: any) {
		return {
			categories: [],
			error: error.message,
		};
	}
};

export const actions: Actions = {
	create: async ({ request, fetch: eventFetch }) => {
		const formData = await request.formData();
		const name = formData.get('name') as string;

		if (!name) {
			return fail(400, { type: 'create', error: 'Nama kategori tidak boleh kosong.' });
		}

		try {
			await createCategory({ name }, eventFetch);
			return { success: true, message: 'Kategori baru berhasil dibuat.' };
		} catch (error: any) {
			return fail(400, { type: 'create', name, error: error.message });
		}
	},

	update: async ({ request, fetch: eventFetch }) => {
		const formData = await request.formData();
		const id = Number(formData.get('id'));
		const name = formData.get('name') as string;

		if (!name) {
			return fail(400, { type: 'update', id, error: 'Nama kategori tidak boleh kosong.' });
		}

		try {
			await updateCategory(id, { name }, eventFetch);
			return { success: true, message: 'Kategori berhasil diperbarui.' };
		} catch (error: any) {
			return fail(400, { type: 'update', id, error: error.message });
		}
	},

	delete: async ({ request, fetch: eventFetch }) => {
		const formData = await request.formData();
		const id = Number(formData.get('id'));

		try {
			await deleteCategory(id, eventFetch);
			return { success: true, message: 'Kategori berhasil dihapus.' };
		} catch (error: any) {
			return fail(400, { type: 'delete', id, error: error.message });
		}
	},
};
