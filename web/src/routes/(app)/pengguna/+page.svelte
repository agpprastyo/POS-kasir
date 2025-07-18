<script lang="ts">
	import type { PageData } from './$types';
	import { goto, invalidateAll } from '$app/navigation';
	import { userProfile } from '$lib/stores';
	import { createUser } from '$lib/api/pengguna';
	import type { CreateUserRequest } from '$lib/types';

	export let data: PageData;

	// --- State untuk Filter ---
	let search = data.queryParams?.search || '';
	let role = data.queryParams?.role || '';
	let isActive = data.queryParams?.is_active === undefined ? '' : String(data.queryParams.is_active);

	// --- State untuk Modal Tambah Pengguna ---
	let createModalDialog: HTMLDialogElement;
	let newUser: CreateUserRequest = {
		username: '',
		email: '',
		password: '',
		role: 'cashier',
		is_active: true,
	};
	let isCreating = false;
	let createError = '';

	// --- Fungsi-fungsi ---
	function formatDate(dateString: string) {
		if (!dateString) return '-';
		return new Date(dateString).toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
		});
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
		const queryString = createQueryString({
			search,
			role,
			is_active: isActive
		});
		goto(`?${queryString}`, { keepFocus: true, noScroll: true });
	}


	let showCreateModal = false;

	function openCreateModal() {
		newUser = { username: '', email: '', password: '', role: 'cashier', is_active: true };
		createError = '';
		showCreateModal = true;
	}
	function closeDialog() {
		showCreateModal = false;
	}

	async function handleCreateUser() {
		isCreating = true;
		createError = '';
		try {
			await createUser(newUser);
			closeDialog();
			await invalidateAll();
		} catch (error: any) {
			createError = error.message;
		} finally {
			isCreating = false;
		}
	}
</script>

<div class="container mx-auto space-y-6">
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold text-gray-800">Manajemen Pengguna</h1>
		{#if $userProfile?.role === 'admin'}
			<button on:click={openCreateModal} class="rounded-md bg-indigo-600 px-4 py-2 text-white shadow-sm hover:bg-indigo-700">
				+ Tambah Pengguna
			</button>
		{/if}
	</div>

	<!-- Filter dan Pencarian -->
	<div class="rounded-lg bg-white p-4 shadow-md">
		<form on:submit|preventDefault={handleFilterSubmit} class="grid grid-cols-1 gap-4 md:grid-cols-4">
			<input type="search" name="search" placeholder="Cari username atau email..." bind:value={search} class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50" />
			<select name="role" bind:value={role} class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50">
				<option value="">Semua Peran</option>
				<option value="admin">Admin</option>
				<option value="manager">Manager</option>
				<option value="cashier">Cashier</option>
			</select>
			<select name="is_active" bind:value={isActive} class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50">
				<option value="">Semua Status</option>
				<option value="true">Aktif</option>
				<option value="false">Tidak Aktif</option>
			</select>
			<button type="submit" class="rounded-md bg-indigo-600 px-4 py-2 text-white shadow-sm hover:bg-indigo-700">
				Terapkan Filter
			</button>
		</form>
	</div>

	<!-- Tabel Pengguna -->
	<div class="overflow-x-auto rounded-lg bg-white shadow-md">
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
			<tr>
				<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Pengguna</th>
				<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Peran</th>
				<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Status</th>
				<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Tanggal Dibuat</th>
				<th class="relative px-6 py-3"><span class="sr-only">Aksi</span></th>
			</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 bg-white">
			{#if data.users && data.users.length > 0}
				{#each data.users as user (user.id)}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap">
							<div class="flex items-center">
								<div class="h-10 w-10 flex-shrink-0">
									<img class="h-10 w-10 rounded-full object-cover" src={user.avatar || `https://ui-avatars.com/api/?name=${user.username}&background=random`} alt={user.username} />
								</div>
								<div class="ml-4">
									<div class="text-sm font-medium text-gray-900">{user.username}</div>
									<div class="text-sm text-gray-500">{user.email}</div>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
                <span class="rounded-full px-2 py-1 text-xs font-semibold capitalize leading-5"
											class:bg-red-100={user.role === 'admin'} class:text-red-800={user.role === 'admin'}
											class:bg-yellow-100={user.role === 'manager'} class:text-yellow-800={user.role === 'manager'}
											class:bg-blue-100={user.role === 'cashier'} class:text-blue-800={user.role === 'cashier'}
								>
                  {user.role}
                </span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							{#if user.is_active}
								<span class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">Aktif</span>
							{:else}
								<span class="inline-flex rounded-full bg-red-100 px-2 text-xs font-semibold leading-5 text-red-800">Tidak Aktif</span>
							{/if}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
							{formatDate(user.created_at)}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
							<a href="/pengguna/{user.id}" class="text-indigo-600 hover:text-indigo-900">Edit</a>
						</td>
					</tr>
				{/each}
			{:else}
				<tr>
					<td colspan="5" class="px-6 py-12 text-center text-gray-500">
						{#if data.error}
							<p class="text-red-500">Error: {data.error}</p>
						{:else}
							<p>Tidak ada data pengguna yang ditemukan.</p>
						{/if}
					</td>
				</tr>
			{/if}
			</tbody>
		</table>
	</div>

	<!-- Paginasi -->
	{#if data.pagination && data.pagination.total_page > 1}
		<div class="flex items-center justify-between rounded-lg bg-white px-4 py-3 shadow-md sm:px-6">
			<div class="text-sm text-gray-700">
				Menampilkan <span class="font-medium">{(data.pagination.current_page - 1) * data.pagination.per_page + 1}</span>
				- <span class="font-medium">{Math.min(data.pagination.current_page * data.pagination.per_page, data.pagination.total_data)}</span>
				dari <span class="font-medium">{data.pagination.total_data}</span> hasil
			</div>
			<div class="flex items-center space-x-2">
				<a href="?{createQueryString({ ...data.queryParams, page: data.pagination.current_page - 1 })}"
					 class="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 {data.pagination.current_page <= 1 ? 'disabled' : ''}"
					 aria-disabled={data.pagination.current_page <= 1}
				>
					Sebelumnya
				</a>
				<a href="?{createQueryString({ ...data.queryParams, page: data.pagination.current_page + 1 })}"
					 class="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 {data.pagination.current_page >= data.pagination.total_page ? 'disabled' : ''}"
					 aria-disabled={data.pagination.current_page >= data.pagination.total_page}
				>
					Berikutnya
				</a>
			</div>
		</div>
	{/if}
</div>

<!-- Modal Tambah Pengguna (Menggunakan <dialog>) -->
{#if showCreateModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center">
		<div class="fixed inset-0 bg-gray-500/75" aria-hidden="true"></div>
		<div class="relative z-10 w-full max-w-lg mx-auto">
			<div class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:p-6">
				<h2 class="mb-4 text-2xl font-bold text-center">Tambah Pengguna Baru</h2>
				<form on:submit|preventDefault={handleCreateUser} class="space-y-4">
					{#if createError}
						<div class="rounded-md bg-red-100 p-3 text-center text-sm text-red-700">{createError}</div>
					{/if}
					<div>
						<label for="username" class="mb-1 block text-sm font-medium text-gray-700">Username</label>
						<input type="text" id="username" bind:value={newUser.username} class="w-full rounded-md border-gray-300" required />
					</div>
					<div>
						<label for="email" class="mb-1 block text-sm font-medium text-gray-700">Email</label>
						<input type="email" id="email" bind:value={newUser.email} class="w-full rounded-md border-gray-300" required />
					</div>
					<div>
						<label for="password" class="mb-1 block text-sm font-medium text-gray-700">Password</label>
						<input type="password" id="password" bind:value={newUser.password} class="w-full rounded-md border-gray-300" required />
					</div>
					<div>
						<label for="role-modal" class="mb-1 block text-sm font-medium text-gray-700">Peran</label>
						<select id="role-modal" bind:value={newUser.role} class="w-full rounded-md border-gray-300">
							<option value="cashier">Cashier</option>
							<option value="manager">Manager</option>
							<option value="admin">Admin</option>
						</select>
					</div>
					<div class="flex items-center">
						<input type="checkbox" id="is_active-modal" bind:checked={newUser.is_active} class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
						<label for="is_active-modal" class="ml-2 block text-sm text-gray-900">Akun Aktif</label>
					</div>
					<div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
						<button type="submit" disabled={isCreating} class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 disabled:cursor-not-allowed disabled:bg-indigo-400 sm:col-start-2">
							{isCreating ? 'Menyimpan...' : 'Simpan Pengguna'}
						</button>
						<button type="button" on:click={closeDialog} class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 shadow-xs ring-gray-300 ring-inset hover:bg-gray-50 sm:col-start-1 sm:mt-0">
							Batal
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<style>
    dialog[max-w-lg] {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        margin: 0;
        padding: 0;
        border: none;
        max-width: 32rem;
        width: 100%;
        z-index: 50;
    }
    .disabled {
        pointer-events: none;
        opacity: 0.6;
    }
</style>
