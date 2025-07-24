import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getProductById, getCategories, updateProduct, updateProductOption, uploadProductImage, uploadProductOptionImage } from '$lib/api/product';

export const load: PageServerLoad = async (event) => {
	const { params, fetch: eventFetch, parent } = event;
	await parent(); // Memastikan otentikasi

	try {
		const [productResponse, categoriesResponse] = await Promise.all([
			getProductById(params.id, eventFetch),
			getCategories(eventFetch)
		]);

		return {
			product: productResponse.data,
			categories: categoriesResponse.data,
		};
	} catch (error) {
		console.error('Gagal memuat detail produk:', error);
		throw redirect(303, '/product');
	}
};

export const actions: Actions = {
	updateDetails: async ({ request, params, fetch: eventFetch }) => {
		const formData = await request.formData();
		const productData = {
			name: formData.get('name') as string,
			category_id: Number(formData.get('category_id')),
			price: Number(formData.get('price')),
			stock: Number(formData.get('stock')),
		};

		try {
			await updateProduct(params.id, productData, eventFetch);
			return { success: true, message: 'Detail produk berhasil diperbarui.' };
		} catch (error: any) {
			return fail(400, { type: 'details', error: error.message });
		}
	},

	updateOption: async ({ request, params, fetch: eventFetch }) => {
		const formData = await request.formData();
		const optionId = formData.get('optionId') as string;
		const optionData = {
			name: formData.get('name') as string,
			additional_price: Number(formData.get('additional_price')),
		};

		try {
			await updateProductOption(params.id, optionId, optionData, eventFetch);
			return { success: true, message: `Opsi "${optionData.name}" berhasil diperbarui.` };
		} catch (error: any) {
			return fail(400, { type: 'option', optionId, error: error.message });
		}
	},

	updateMainImage: async ({ request, params, fetch: eventFetch }) => {
		const formData = await request.formData();
		const image = formData.get('image') as File;

		if (!image || image.size === 0) {
			return fail(400, { type: 'mainImage', error: 'Silakan pilih file gambar.' });
		}

		try {
			await uploadProductImage(params.id, image, eventFetch);
			return { success: true, message: 'Gambar produk utama berhasil diunggah.' };
		} catch (error: any) {
			return fail(400, { type: 'mainImage', error: error.message });
		}
	},

	updateOptionImage: async ({ request, params, fetch: eventFetch }) => {
		const formData = await request.formData();
		const image = formData.get('image') as File;
		const optionId = formData.get('optionId') as string;

		if (!image || image.size === 0) {
			return fail(400, { type: 'optionImage', optionId, error: 'Silakan pilih file gambar.' });
		}

		try {
			await uploadProductOptionImage(params.id, optionId, image, eventFetch);
			return { success: true, message: 'Gambar opsi berhasil diunggah.' };
		} catch (error: any) {
			return fail(400, { type: 'optionImage', optionId, error: error.message });
		}
	}
};
