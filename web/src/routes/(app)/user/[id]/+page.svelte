<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';
	import type { PageData, ActionData } from './$types';

	export let data: PageData;
	export let form: ActionData;

	// State untuk mengelola pesan sukses atau error dari form actions
	let message: { type: 'success' | 'error'; text: string } | null = null;
	$: {
		if (form?.success) {
			message = { type: 'success', text: form.message };
		} else if (form?.error) {
			message = { type: 'error', text: form.error };
		}
	}

	// State untuk modal konfirmasi hapus
	let deleteModal: HTMLDialogElement;
</script>

<div class="container mx-auto max-w-4xl space-y-8">
	<div class="flex items-center gap-4">
		<a href="/user" class="rounded-full p-2 hover:bg-gray-100">
			<svg class="h-6 w-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
		</a>
		<h1 class="text-3xl font-bold text-gray-800">Edit Pengguna: {data.user.username}</h1>
	</div>

	<!-- Notifikasi Sukses/Error -->
	{#if message}
		<div class="rounded-md p-4 text-center" class:bg-green-100={message.type === 'success'} class:text-green-800={message.type === 'success'} class:bg-red-100={message.type === 'error'} class:text-red-800={message.type === 'error'}>
			{message.text}
		</div>
	{/if}

	<!-- Form Edit Data Pengguna -->
	<div class="rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-semibold">Informasi Pengguna</h2>
		<form method="POST" action="?/update" use:enhance>
			<div class="space-y-4">
				<div>
					<label for="username" class="mb-1 block text-sm font-medium text-gray-700">Username</label>
					<input type="text" id="username" name="username" value={data.user.username} class="w-full rounded-md border-gray-300" required />
				</div>
				<div>
					<label for="email" class="mb-1 block text-sm font-medium text-gray-700">Email</label>
					<input type="email" id="email" name="email" value={data.user.email} class="w-full rounded-md border-gray-300" required />
				</div>
				<div>
					<label for="role" class="mb-1 block text-sm font-medium text-gray-700">Peran</label>
					<select id="role" name="role" class="w-full rounded-md border-gray-300">
						<option value="cashier" selected={data.user.role === 'cashier'}>Cashier</option>
						<option value="manager" selected={data.user.role === 'manager'}>Manager</option>
						<option value="admin" selected={data.user.role === 'admin'}>Admin</option>
					</select>
				</div>
				<button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-white shadow-sm hover:bg-indigo-700">Simpan Perubahan</button>
			</div>
		</form>
	</div>

	<!-- Aksi Berbahaya -->
	<div class="rounded-lg border-2 border-red-300 bg-red-50 p-6 shadow-md">
		<h2 class="mb-4 text-xl font-semibold text-red-800">Zona Berbahaya</h2>
		<div class="flex flex-wrap items-center justify-between gap-4">
			<div>
				<p class="font-medium">Ubah Status Akun</p>
				<p class="text-sm text-gray-600">
					Saat ini akun ini {data.user.is_active ? 'aktif' : 'tidak aktif'}. Menekan tombol ini akan membalikkan statusnya.
				</p>
			</div>
			<form method="POST" action="?/toggleStatus" use:enhance>
				<button type="submit" class="rounded-md bg-yellow-500 px-4 py-2 text-white shadow-sm hover:bg-yellow-600">
					{data.user.is_active ? 'Nonaktifkan Akun' : 'Aktifkan Akun'}
				</button>
			</form>
		</div>
		<hr class="my-6 border-red-200" />
		<div class="flex flex-wrap items-center justify-between gap-4">
			<div>
				<p class="font-medium">Hapus Pengguna</p>
				<p class="text-sm text-gray-600">Aksi ini tidak dapat diurungkan. Semua data pengguna akan dihapus permanen.</p>
			</div>
			<button type="button" on:click={() => deleteModal.showModal()} class="rounded-md bg-red-600 px-4 py-2 text-white shadow-sm hover:bg-red-700">
				Hapus Pengguna Ini
			</button>
		</div>
	</div>
</div>

<!-- Modal Konfirmasi Hapus -->
<dialog bind:this={deleteModal} class="m-auto max-w-md rounded-lg p-0 shadow-xl backdrop:bg-black/50">
	<div class="p-6 text-center">
		<div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-red-100">
			<svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
		</div>
		<h3 class="mt-4 text-lg font-medium text-gray-900">Hapus Pengguna?</h3>
		<p class="mt-2 text-sm text-gray-500">
			Apakah Anda yakin ingin menghapus pengguna "{data.user.username}"? Aksi ini tidak dapat dibatalkan.
		</p>
		<div class="mt-6 flex justify-center gap-4">
			<button type="button" on:click={() => deleteModal.close()} class="rounded-md bg-gray-200 px-4 py-2 text-gray-800 hover:bg-gray-300">Batal</button>
			<form method="POST" action="?/delete" use:enhance>
				<button type="submit" class="rounded-md bg-red-600 px-4 py-2 text-white shadow-sm hover:bg-red-700">Ya, Hapus</button>
			</form>
		</div>
	</div>
</dialog>
