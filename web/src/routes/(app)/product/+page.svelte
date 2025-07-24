<script lang="ts">
	import type { PageData } from './$types';
	import { goto, invalidateAll } from '$app/navigation';
	import { createProduct, uploadProductImage, uploadProductOptionImage } from '$lib/api/product';
	import type { CreateProductRequest, Product } from '$lib/types';
	import { SvelteURLSearchParams } from 'svelte/reactivity';

	export let data: PageData;

	// --- State untuk Filter ---
	let search = data.queryParams?.search || '';
	let category_id = data.queryParams?.category_id || '';

	// --- State untuk Modal Tambah Produk (Satu Langkah) ---
	let showCreateModal = false;
	let newProduct: CreateProductRequest = {
		name: '',
		category_id: 0,
		price: 0,
		stock: 0,
		options: []
	};
	let mainImageFile: FileList | null = null;
	let optionImageFiles: File[] = []; // Array sederhana untuk menampung file opsi

	let isCreating = false;
	let createError = '';

	// --- Fungsi-fungsi ---
	function formatCurrency(value: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value);
	}

	function createQueryString(params: Record<string, any>): string {
		const cleanedParams: Record<string, string> = {};
		for (const key in params) {
			if (params[key] !== null && params[key] !== undefined && params[key] !== '') {
				cleanedParams[key] = String(params[key]);
			}
		}
		return new URLSearchParams(cleanedParams).toString();
	}

	function handleFilterSubmit() {
		const params = new SvelteURLSearchParams();
		if (search) params.set('search', search);
		if (category_id) params.set('category_id', category_id);
		params.set('view', data.view);
		goto(`?${params.toString()}`, { keepFocus: true, noScroll: true });
	}

	function addOption() {
		newProduct.options = [...(newProduct.options || []), { name: '', additional_price: 0 }];
		optionImageFiles = [...optionImageFiles, undefined as any]; // Tambah slot kosong untuk file
	}

	function removeOption(index: number) {
		newProduct.options = (newProduct.options || []).filter((_, i) => i !== index);
		optionImageFiles = optionImageFiles.filter((_, i) => i !== index); // Hapus slot file
	}

	async function handleCreateProduct() {
		isCreating = true;
		createError = '';
		try {
			// Langkah 1: Buat produk dengan detail teks
			const response = await createProduct(newProduct);
			const createdProduct = response.data;

			// Langkah 2: Unggah gambar utama jika ada
			if (mainImageFile && mainImageFile[0]) {
				await uploadProductImage(createdProduct.id, mainImageFile[0]);
			}

			// Langkah 3: Unggah gambar untuk setiap opsi jika ada
			if (createdProduct.options) {
				for (let i = 0; i < createdProduct.options.length; i++) {
					const option = createdProduct.options[i];
					const file = optionImageFiles[i];
					if (option.id && file) {
						await uploadProductOptionImage(createdProduct.id, option.id, file);
					}
				}
			}

			closeModal();
			await invalidateAll(); // Muat ulang data di halaman
		} catch (error: any) {
			createError = error.message;
		} finally {
			isCreating = false;
		}
	}

	function openModal() {
		// Reset state
		newProduct = { name: '', category_id: 0, price: 0, stock: 0, options: [] };
		mainImageFile = null;
		optionImageFiles = [];
		createError = '';
		showCreateModal = true;
	}

	function closeModal() {
		showCreateModal = false;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeModal();
		}
	}
</script>


<svelte:window on:keydown={handleKeydown} />

<div class="container mx-auto space-y-6">
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold text-gray-800">Manajemen Produk</h1>
		<button on:click={openModal} class="rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700 transition-colors">
			+ Tambah Produk
		</button>
	</div>

	<!-- Filter dan Kontrol Tampilan -->
	<div class="rounded-xl bg-white p-4 shadow-sm shadow-gray-200">
		<div class="flex flex-wrap items-center justify-between gap-4">
			<form on:submit|preventDefault={handleFilterSubmit} class="flex flex-grow items-center gap-4">
				<div class="flex flex-grow items-center gap-2 rounded-lg bg-gray-100 px-4 py-2 md:w-auto">
					<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35M17 11a6 6 0 11-12 0 6 6 0 0112 0z"></path></svg>
					<input type="search" name="search" placeholder="Cari produk..." bind:value={search} class="w-full bg-transparent text-gray-700 placeholder-gray-400 outline-none border-0 rounded-lg" />
				</div>
				<select name="category_id" bind:value={category_id} class="rounded-lg border-1 border-gray-300 bg-gray-50 px-4 py-2 text-gray-700 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100">
					<option value="">Semua Kategori</option>
					{#each data.categories as category (category.id)}
						<option value={category.id}>{category.name}</option>
					{/each}
				</select>
				<button type="submit" class="rounded-lg bg-indigo-600 px-6 py-2 font-semibold text-white shadow transition hover:bg-indigo-700">Filter</button>
			</form>

			<!-- Tombol Ganti Tampilan -->
			<div class="flex rounded-lg bg-gray-100 p-1">
				<a href="?{createQueryString({ ...data.queryParams, view: 'grid' })}" class="rounded-md px-3 py-1.5 text-sm font-medium" class:bg-white={data.view === 'grid'} class:text-gray-500={data.view !== 'grid'} class:shadow={data.view === 'grid'}>Grid</a>
				<a href="?{createQueryString({ ...data.queryParams, view: 'list' })}" class="rounded-md px-3 py-1.5 text-sm font-medium" class:bg-white={data.view === 'list'} class:text-gray-500={data.view !== 'list'} class:shadow={data.view === 'list'}>List</a>
			</div>
		</div>
	</div>

	<!-- Tampilan Produk (Grid atau List) -->
	{#if data.view === 'grid'}
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
			{#each data.products as product (product.id)}

				<div class="group rounded-lg bg-white shadow-sm shadow-gray-200 overflow-hidden transition-shadow hover:shadow-md">
					<div class="relative w-full overflow-hidden bg-gray-200" style="aspect-ratio:1/1;">
						<img src={product.image_url || `https://placehold.co/400x400/e2e8f0/64748b?text=${product.name}`} alt={product.name} class="absolute inset-0 h-full w-full object-cover object-center transition-transform group-hover:scale-105" />
					</div>
					<div class="p-4">
						<h3 class="text-sm text-gray-500">{product.category_name}</h3>
						<a href="/product/{product.id}">
							<p class="mt-1 font-semibold text-gray-900 hover:text-indigo-600">{product.name}</p>
						</a>
						<p class="mt-2 text-lg font-bold text-gray-800">{formatCurrency(product.price)}</p>
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<div class="overflow-x-auto rounded-lg bg-white shadow-md">
			<table class="min-w-full divide-y divide-gray-200">
				<thead class="bg-gray-50">
				<tr>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Produk</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">SKU</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Kategori</th>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Harga</th>
					<th class="relative px-6 py-3"><span class="sr-only">Aksi</span></th>
				</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 bg-white">
				{#each data.products as product (product.id)}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap">
							<div class="flex items-center">
								<div class="h-10 w-10 flex-shrink-0">
									<img class="h-10 w-10 rounded-md object-cover" src={product.image_url || `https://placehold.co/100x100/e2e8f0/64748b?text=Img`} alt={product.name} />
								</div>
								<div class="ml-4">
									<div class="text-sm font-medium text-gray-900">{product.name}</div>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{product.sku || '-'}</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{product.category_name}</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-800">{formatCurrency(product.price)}</td>
						<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
							<a href="/product/{product.id}" class="text-indigo-600 hover:text-indigo-900">Detail</a>
						</td>
					</tr>
				{/each}
				</tbody>
			</table>
		</div>
	{/if}

	<!-- Pesan jika tidak ada produk -->
	{#if data.products.length === 0 || data.products == null}
		<div class="rounded-lg bg-white p-12 text-center text-gray-500 shadow-md">
			{#if data.error}
				<p class="text-red-500">Error: {data.error}</p>
			{:else}
				<p>Tidak ada produk yang ditemukan.</p>
			{/if}
		</div>
	{/if}

	<!-- Paginasi -->
	{#if data.pagination && data.pagination.total_page > 1}
		<div class="flex items-center justify-between rounded-lg bg-white px-4 py-3 shadow-lg sm:px-6">
			<div class="text-sm text-gray-700">
				Menampilkan <span class="font-medium">{(data.pagination.current_page - 1) * data.pagination.per_page + 1}</span>
				- <span class="font-medium">{Math.min(data.pagination.current_page * data.pagination.per_page, data.pagination.total_data)}</span>
				dari <span class="font-medium">{data.pagination.total_data}</span> hasil
			</div>
			<div class="flex items-center space-x-2">
				<a href="?{createQueryString({ ...data.queryParams, page: data.pagination.current_page - 1, view: data.view })}"
					 class="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 {data.pagination.current_page <= 1 ? 'disabled' : ''}"
					 aria-disabled={data.pagination.current_page <= 1}
				>
					Sebelumnya
				</a>
				<a href="?{createQueryString({ ...data.queryParams, page: data.pagination.current_page + 1, view: data.view })}"
					 class="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 {data.pagination.current_page >= data.pagination.total_page ? 'disabled' : ''}"
					 aria-disabled={data.pagination.current_page >= data.pagination.total_page}
				>
					Berikutnya
				</a>
			</div>
		</div>
	{/if}
</div>

<!-- Modal Tambah Produk (Satu Langkah) -->
{#if showCreateModal}
	<dialog aria-modal="true" tabindex="0" aria-labelledby="dialogTitle" class="fixed inset-0 z-50 flex items-center justify-center w-full h-full bg-transparent backdrop-blur-sm" on:click|self={closeModal}>
		<div class="w-full max-w-2xl rounded-lg bg-white shadow-xl max-h-[90vh] overflow-y-auto">
			<form on:submit|preventDefault={handleCreateProduct} class="p-6">
				<h2 class="mb-6 text-2xl font-bold">Tambah Produk Baru</h2>
				{#if createError}<p class="mb-4 text-sm text-red-600">{createError}</p>{/if}

				<div class="space-y-4">
					<!-- Detail Produk Utama -->
					<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
						<div>
							<label for="name" class="block text-sm font-medium text-gray-700">Nama Produk</label>
							<input type="text" bind:value={newProduct.name} id="name" class="mt-1 block w-full rounded-md border-gray-300" required />
						</div>
						<div>
							<label for="category" class="block text-sm font-medium text-gray-700">Kategori</label>
							<select bind:value={newProduct.category_id} id="category" class="mt-1 block w-full rounded-md border-gray-300" required>
								<option value={0} disabled>Pilih kategori</option>
								{#each data.categories as category (category.id)}<option value={category.id}>{category.name}</option>{/each}
							</select>
						</div>
						<div>
							<label for="price" class="block text-sm font-medium text-gray-700">Harga</label>
							<input type="number" bind:value={newProduct.price} id="price" class="mt-1 block w-full rounded-md border-gray-300" required />
						</div>
						<div>
							<label for="stock" class="block text-sm font-medium text-gray-700">Stok</label>
							<input type="number" bind:value={newProduct.stock} id="stock" class="mt-1 block w-full rounded-md border-gray-300" required />
						</div>
					</div>
					<div>
						<label for="main-image-upload" class="block text-sm font-medium text-gray-700">Gambar Produk Utama (Opsional)</label>
						<input type="file" id="main-image-upload" accept="image/*" bind:files={mainImageFile} class="mt-1 block w-full text-sm" />
					</div>

					<hr class="my-4" />

					<!-- Opsi Produk -->
					<h3 class="text-lg font-semibold">Opsi Produk</h3>
					<div class="space-y-4">
						{#if newProduct.options}
							{#each newProduct.options as option, i (i)}
								<div class="rounded-md border border-gray-200 p-4">
									<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
										<input type="text" placeholder="Nama Opsi (e.g., Large)" bind:value={option.name} class="rounded-md border-gray-300" />
										<input type="number" placeholder="Harga Tambahan" bind:value={option.additional_price} class="rounded-md border-gray-300" />
									</div>
									<div class="mt-2">
										<label for="option-image-{i}" class="text-sm font-medium text-gray-700">Gambar Opsi (Opsional)</label>
										<input type="file" id="option-image-{i}" accept="image/*" on:change={(e) => optionImageFiles[i] = e.currentTarget.files?.[0]} class="mt-1 block w-full text-sm" />
									</div>
									<button type="button" on:click={() => removeOption(i)} class="mt-2 text-xs text-red-500 hover:text-red-700">Hapus Opsi</button>
								</div>
							{/each}
						{/if}
					</div>
					<button type="button" on:click={addOption} class="mt-4 text-sm text-indigo-600 hover:text-indigo-800">+ Tambah Opsi</button>
				</div>

				<!-- Tombol Aksi -->
				<div class="mt-8 flex justify-end gap-4 border-t border-gray-200 pt-4">
					<button type="button" on:click={closeModal} class="rounded-md bg-gray-200 px-4 py-2">Batal</button>
					<button type="submit" disabled={isCreating} class="rounded-md bg-indigo-600 px-4 py-2 text-white disabled:bg-indigo-400">
						{isCreating ? 'Menyimpan...' : 'Simpan Produk'}
					</button>
				</div>
			</form>
		</div>
	</dialog>
{/if}

<style>
    .disabled {
        pointer-events: none;
        opacity: 0.6;
    }
</style>
