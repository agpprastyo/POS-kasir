<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData, ActionData } from './$types';
	import type { CategoryWithCount } from '$lib/types';

	export let data: PageData;
	export let form: ActionData;

	// State untuk modal edit dan delete
	let showEditModal = false;
	let showDeleteModal = false;
	let selectedCategory: CategoryWithCount | null = null;

	function openEditModal(category: CategoryWithCount) {
		selectedCategory = category;
		showEditModal = true;
	}

	function openDeleteModal(category: CategoryWithCount) {
		selectedCategory = category;
		showDeleteModal = true;
	}

	// Tutup modal setelah aksi berhasil
	$: if (form?.success) {
		showEditModal = false;
		showDeleteModal = false;
	}

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleString('id-ID', { dateStyle: 'long', timeStyle: 'short' });
	}
</script>

<div class="container mx-auto  space-y-8">
	<h1 class="text-3xl font-bold text-gray-800">Manajemen Kategori</h1>

	{#if form?.success}
		<div class="rounded-lg bg-green-50 p-4 text-center text-sm font-medium text-green-700 border border-green-200">{form.message}</div>
	{/if}
	{#if form?.error && !form.id}
		<div class="rounded-lg bg-red-50 p-4 text-center text-sm font-medium text-red-700 border border-red-200">{form.error}</div>
	{/if}

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
		<div class="lg:col-span-2">
			<div class="overflow-x-auto rounded-xl bg-white border border-gray-200">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Nama Kategori</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Jumlah Produk</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Terakhir Diperbarui</th>
						<th class="relative px-6 py-3"><span class="sr-only">Aksi</span></th>
					</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 bg-white">
					{#if data.categories && data.categories.length > 0}
						{#each data.categories as category (category.id)}
							<tr>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{category.name}</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{category.product_count}</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(category.updated_at)}</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-4">
									<button on:click={() => openEditModal(category)} class="text-indigo-600 hover:text-indigo-900">Edit</button>
									<button on:click={() => openDeleteModal(category)} class="text-red-600 hover:text-red-900">Hapus</button>
								</td>
							</tr>
						{/each}
					{:else}
						<tr>
							<td colspan="4" class="px-6 py-12 text-center text-gray-500">
								{#if data.error}<p class="text-red-500">{data.error}</p>{:else}<p>Belum ada kategori.</p>{/if}
							</td>
						</tr>
					{/if}
					</tbody>
				</table>
			</div>
		</div>

		<div class="lg:col-span-1">
			<div class="rounded-xl bg-white p-6 shadow-sm border border-gray-200 sticky top-6">
				<h2 class="text-xl font-semibold text-gray-800 mb-4">Tambah Kategori Baru</h2>
				<form method="POST" action="?/create" use:enhance class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-600">Nama Kategori</label>
						<input type="text" name="name" id="name" class="mt-1 w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
					</div>
					<button type="submit" class="w-full rounded-lg bg-indigo-600 px-4 py-2.5 font-semibold text-white shadow-sm hover:bg-indigo-700 transition-colors">Simpan Kategori</button>
				</form>
			</div>
		</div>
	</div>
</div>

{#if showEditModal && selectedCategory}
	<dialog open class="fixed inset-0 z-50 m-auto max-w-md rounded-lg p-0 shadow-xl backdrop:bg-black/50">
		<div class="p-6">
			<h3 class="text-lg font-medium text-gray-900">Edit Kategori</h3>
			{#if form?.type === 'update' && form?.id === selectedCategory.id && form?.error}<p class="mt-2 text-sm text-red-600">{form.error}</p>{/if}
			<form method="POST" action="?/update" use:enhance class="mt-4 space-y-4">
				<input type="hidden" name="id" value={selectedCategory.id} />
				<div>
					<label for="edit-name" class="block text-sm font-medium">Nama Kategori</label>
					<input type="text" name="name" id="edit-name" value={selectedCategory.name} class="mt-1 w-full rounded-md border-gray-300" />
				</div>
				<div class="mt-6 flex justify-end gap-4">
					<button type="button" on:click={() => showEditModal = false} class="rounded-md bg-gray-200 px-4 py-2 text-gray-800 hover:bg-gray-300">Batal</button>
					<button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-white">Simpan</button>
				</div>
			</form>
		</div>
	</dialog>
{/if}

{#if showDeleteModal && selectedCategory}
	<dialog open class="fixed inset-0 z-50 m-auto max-w-md rounded-lg p-0 shadow-xl backdrop:bg-black/50">
		<div class="p-6 text-center">
			<h3 class="mt-4 text-lg font-medium text-gray-900">Hapus Kategori?</h3>
			<p class="mt-2 text-sm text-gray-500">
				Yakin ingin menghapus kategori "{selectedCategory.name}"? Aksi ini tidak dapat dibatalkan.
			</p>
			{#if form?.type === 'delete' && form?.id === selectedCategory.id && form?.error}<p class="mt-2 text-sm text-red-600">{form.error}</p>{/if}
			<div class="mt-6 flex justify-center gap-4">
				<button type="button" on:click={() => showDeleteModal = false} class="rounded-md bg-gray-200 px-4 py-2 text-gray-800 hover:bg-gray-300">Batal</button>
				<form method="POST" action="?/delete" use:enhance>
					<input type="hidden" name="id" value={selectedCategory.id} />
					<button type="submit" class="rounded-md bg-red-600 px-4 py-2 text-white">Ya, Hapus</button>
				</form>
			</div>
		</div>
	</dialog>
{/if}
