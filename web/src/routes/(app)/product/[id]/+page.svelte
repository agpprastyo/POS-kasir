<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData, ActionData } from './$types';
	export let data: PageData;
	export let form: ActionData;

</script>

<div class="container mx-auto  space-y-8">
	<div class="flex items-center gap-4">
		<a href="/product" aria-label="product" class="rounded-full p-2 hover:bg-gray-100 transition-colors">
			<svg class="h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
		</a>
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Kembali ke Daftar Produk</h1>
			<p class="text-sm text-gray-500">Kelola detail produk dan variasinya di bawah ini.</p>
		</div>
	</div>

	<!-- Notifikasi Global -->
	{#if form?.success}
		<div class="rounded-lg bg-green-50 p-4 text-center text-sm font-medium text-green-700 border border-green-200">{form.message}</div>
	{/if}

	<!-- Layout Utama -->
	<div class="grid grid-cols-1 lg:grid-cols-5 gap-8">

		<!-- Kolom Kiri: Form Edit Detail -->
		<div class="lg:col-span-3">
			<div class="rounded-xl bg-white p-6 shadow-sm border border-gray-200">
				<h2 class="text-xl font-semibold text-gray-800 border-b pb-4 mb-6">Informasi Produk</h2>
				{#if form?.type === 'details' && form?.error}<p class="mb-4 text-sm text-red-600">{form.error}</p>{/if}
				<form method="POST" action="?/updateDetails" use:enhance class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-600">Nama Produk</label>
						<input type="text" name="name" value={data.product.name} class="mt-1 w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
					</div>
					<div>
						<label for="category_id" class="block text-sm font-medium text-gray-600">Kategori</label>
						<select name="category_id" class="mt-1 w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
							{#each data.categories as category (category.id)}
								<option value={category.id} selected={category.id === data.product.category_id}>{category.name}</option>
							{/each}
						</select>
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="price" class="block text-sm font-medium text-gray-600">Harga</label>
							<input type="number" name="price" value={data.product.price} class="mt-1 w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
						</div>
						<div>
							<label for="stock" class="block text-sm font-medium text-gray-600">Stok</label>
							<input type="number" name="stock" value={data.product.stock} class="mt-1 w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
						</div>
					</div>
					<button type="submit" class="w-full rounded-lg bg-indigo-600 px-4 py-2.5 font-semibold text-white shadow-sm hover:bg-indigo-700 transition-all focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">Simpan Perubahan</button>
				</form>
			</div>
		</div>

		<!-- Kolom Kanan: Gambar Utama -->

		<div class="lg:col-span-2">
			<div class="rounded-xl bg-white p-6 shadow-sm border border-gray-200 sticky top-6">
				<h3 class="text-xl font-semibold text-gray-800 mb-4">Gambar Utama</h3>
				<img src={data.product.image_url || `https://placehold.co/800x800/f3f4f6/9ca3af?text=${data.product.name}`} alt={data.product.name} class="aspect-square w-full rounded-lg object-cover" />
				<form method="POST" action="?/updateMainImage" use:enhance enctype="multipart/form-data" class="mt-4">
					<label for="main-image" class="text-sm font-medium text-gray-600">Ubah Gambar</label>
					<div class="mt-2 flex gap-2">
						<input type="file" name="image" id="main-image" class="block w-full text-sm text-gray-500 file:mr-4 file:rounded-lg file:border-0 file:bg-gray-100 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-indigo-600 hover:file:bg-indigo-100" />
						<button type="submit" class="flex-shrink-0 rounded-lg bg-gray-700 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-gray-800 transition-colors">Unggah</button>
					</div>
					{#if form?.type === 'mainImage' && form?.error}<p class="mt-1 text-sm text-red-600">{form.error}</p>{/if}
				</form>
			</div>
		</div>

		<!-- Opsi Produk (di bawah, span penuh) -->
		<div class="lg:col-span-5 rounded-xl bg-white p-6 shadow-sm border border-gray-200">
			<h2 class="text-xl font-semibold text-gray-800 border-b pb-4 mb-6">Opsi Produk</h2>
			<div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-6">
				{#if data.product.options && data.product.options.length > 0}
					{#each data.product.options as option (option.id)}
						<div class="rounded-lg border border-gray-200 bg-gray-50 p-4 space-y-4">
							<img src={option.image_url || 'https://placehold.co/400x400/f3f4f6/9ca3af?text=Opsi'} alt={option.name} class="aspect-square w-full rounded-md object-cover" />

							<!-- Form Edit Opsi -->
							<form method="POST" action="?/updateOption" use:enhance class="space-y-3">
								{#if form?.type === 'option' && form?.optionId === option.id && form?.error}<p class="mb-2 text-sm text-red-600">{form.error}</p>{/if}
								<input type="hidden" name="optionId" value={option.id} />
								<div>
									<label for="option-name-{option.id}" class="block text-xs font-medium text-gray-500">Nama Opsi</label>
									<input type="text" name="name" id="option-name-{option.id}" value={option.name} class="mt-1 w-full rounded-md border-gray-300 text-sm shadow-sm" />
								</div>
								<div>
									<label for="option-price-{option.id}" class="block text-xs font-medium text-gray-500">Harga Tambahan</label>
									<input type="number" name="additional_price" id="option-price-{option.id}" value={option.additional_price} class="mt-1 w-full rounded-md border-gray-300 text-sm shadow-sm" />
								</div>
								<button type="submit" class="w-full rounded-md bg-indigo-500 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-600 transition-colors">Simpan Opsi</button>
							</form>

							<!-- Form Upload Gambar Opsi -->
							<form method="POST" action="?/updateOptionImage" use:enhance enctype="multipart/form-data">
								<input type="hidden" name="optionId" value={option.id} />
								<label for="option-image-{option.id}" class="text-xs font-medium text-gray-500">Ubah Gambar Opsi</label>
								<div class="mt-1 flex gap-2">
									<input type="file" name="image" id="option-image-{option.id}" class="block w-full text-xs text-gray-500 file:mr-2 file:rounded file:border-0 file:bg-white file:px-2 file:py-1 file:text-xs file:font-semibold" />
									<button type="submit" class="flex-shrink-0 rounded bg-gray-600 px-2 py-1 text-xs text-white hover:bg-gray-700 transition-colors">Unggah</button>
								</div>
								{#if form?.type === 'optionImage' && form?.optionId === option.id && form?.error}<p class="mt-1 text-xs text-red-600">{form.error}</p>{/if}
							</form>

							<!-- Tombol Hapus Opsi -->
							<form method="POST" action="?/deleteOption" use:enhance>
								<input type="hidden" name="optionId" value={option.id} />
								<button type="submit" class="w-full rounded-md bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 shadow-sm hover:bg-red-100 border border-red-200 transition-colors">Hapus Opsi</button>
							</form>
						</div>
					{/each}
				{/if}

				<!-- Form Tambah Opsi Baru -->
				<div class="rounded-lg border-2 border-dashed border-gray-300 bg-white p-4">
					<h3 class="font-semibold text-gray-700">Tambah Opsi Baru</h3>
					<form method="POST" action="?/createOption" use:enhance class="mt-4 space-y-3">
						{#if form?.type === 'createOption' && form?.error}<p class="text-sm text-red-600">{form.error}</p>{/if}
						<div>
							<label for="new-option-name" class="block text-xs font-medium text-gray-500">Nama Opsi</label>
							<input type="text" name="name" id="new-option-name" class="mt-1 w-full rounded-md border-gray-300 text-sm shadow-sm" />
						</div>
						<div>
							<label for="new-option-price" class="block text-xs font-medium text-gray-500">Harga Tambahan</label>
							<input type="number" name="additional_price" id="new-option-price" value="0" class="mt-1 w-full rounded-md border-gray-300 text-sm shadow-sm" />
						</div>
						<button type="submit" class="w-full rounded-md bg-green-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-green-700 transition-colors">Tambah Opsi</button>
					</form>
				</div>
			</div>
		</div>
	</div>
</div>
