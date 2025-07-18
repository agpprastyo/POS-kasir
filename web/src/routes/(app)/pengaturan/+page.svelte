<script lang="ts">
	import { userProfile } from '$lib/stores';
	import { updateAvatar, updatePassword } from '$lib/api/pengaturan';
	import type { Profile } from '$lib/types';

	// State untuk form ubah avatar
	let avatarFile: FileList | null = null;
	let avatarPreview: string | null = null;
	let isAvatarLoading = false;
	let avatarMessage = { type: '', text: '' };

	// State untuk form ubah password
	let oldPassword = '';
	let newPassword = '';
	let isPasswordLoading = false;
	let passwordMessage = { type: '', text: '' };

	// Fungsi untuk menampilkan preview gambar saat dipilih
	function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			const file = input.files[0];
			avatarPreview = URL.createObjectURL(file);
		}
	}

	// Fungsi untuk menangani submit form avatar
	async function handleAvatarUpdate() {
		if (!avatarFile || avatarFile.length === 0) {
			avatarMessage = { type: 'error', text: 'Silakan pilih file gambar terlebih dahulu.' };
			return;
		}

		isAvatarLoading = true;
		avatarMessage = { type: '', text: '' };

		try {
			const result = await updateAvatar(avatarFile[0]);
			// Perbarui userProfile store dengan data baru dari API
			userProfile.set(result.data);
			avatarMessage = { type: 'success', text: 'Avatar berhasil diperbarui!' };
			avatarPreview = null; // Hapus preview setelah berhasil
			(document.getElementById('avatar-form') as HTMLFormElement).reset();
		} catch (error: any) {
			avatarMessage = { type: 'error', text: error.message };
		} finally {
			isAvatarLoading = false;
		}
	}

	// Fungsi untuk menangani submit form password
	async function handlePasswordUpdate() {
		if (!oldPassword || !newPassword) {
			passwordMessage = { type: 'error', text: 'Semua kolom wajib diisi.' };
			return;
		}
		isPasswordLoading = true;
		passwordMessage = { type: '', text: '' };

		try {
			await updatePassword({ old_password: oldPassword, new_password: newPassword });
			passwordMessage = { type: 'success', text: 'Password berhasil diperbarui!' };
			(document.getElementById('password-form') as HTMLFormElement).reset();
			oldPassword = '';
			newPassword = '';
		} catch (error: any) {
			passwordMessage = { type: 'error', text: error.message };
		} finally {
			isPasswordLoading = false;
		}
	}
</script>

<div class="container mx-auto max-w-4xl space-y-8">
	<h1 class="text-3xl font-bold text-gray-800">Pengaturan Akun</h1>

	<!-- Bagian Profil Pengguna -->
	{#if $userProfile}
		<div class="rounded-lg bg-white p-6 shadow-md">
			<h2 class="mb-4 text-xl font-semibold">Profil Anda</h2>
			<div class="flex items-center space-x-6">
				<img
					class="h-24 w-24 rounded-full object-cover ring-4 ring-indigo-200"
					src={avatarPreview || $userProfile.avatar || `https://ui-avatars.com/api/?name=${$userProfile.username}&background=random`}
					alt="Avatar Pengguna"
				/>
				<div>
					<p class="text-2xl font-bold text-gray-900">{$userProfile.username}</p>
					<p class="text-gray-500">{$userProfile.email}</p>
					<p class="mt-1 rounded-full bg-blue-100 px-3 py-1 text-sm font-medium capitalize text-blue-800">
						Peran: {$userProfile.role}
					</p>
				</div>
			</div>
		</div>
	{/if}

	<!-- Form Ubah Avatar -->
	<div class="rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-semibold">Ubah Foto Profil</h2>
		<form id="avatar-form" on:submit|preventDefault={handleAvatarUpdate}>
			{#if avatarMessage.text}
				<div class="mb-4 rounded-md p-3 text-center text-sm" class:bg-green-100={avatarMessage.type === 'success'} class:text-green-700={avatarMessage.type === 'success'} class:bg-red-100={avatarMessage.type === 'error'} class:text-red-700={avatarMessage.type === 'error'} >
					{avatarMessage.text}
				</div>
			{/if}
			<div class="mb-4">
				<label for="avatar" class="mb-2 block text-sm font-medium text-gray-700">Pilih Gambar Baru (JPG/PNG)</label>
				<input
					type="file"
					id="avatar"
					name="avatar"
					accept="image/png, image/jpeg"
					bind:files={avatarFile}
					on:change={handleFileSelect}
					class="block w-full text-sm text-gray-500 file:mr-4 file:rounded-md file:border-0 file:bg-indigo-50 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-indigo-600 hover:file:bg-indigo-100"
				/>
			</div>
			<button type="submit" disabled={isAvatarLoading || !avatarFile} class="rounded-md bg-indigo-600 px-4 py-2 text-white shadow-sm hover:bg-indigo-700 disabled:cursor-not-allowed disabled:bg-indigo-400">
				{isAvatarLoading ? 'Mengunggah...' : 'Simpan Avatar'}
			</button>
		</form>
	</div>

	<!-- Form Ubah Password -->
	<div class="rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-semibold">Ubah Password</h2>
		<form id="password-form" on:submit|preventDefault={handlePasswordUpdate} class="space-y-4">
			{#if passwordMessage.text}
				<div class="rounded-md p-3 text-center text-sm" class:bg-green-100={passwordMessage.type === 'success'} class:text-green-700={passwordMessage.type === 'success'} class:bg-red-100={passwordMessage.type === 'error'} class:text-red-700={passwordMessage.type === 'error'} >
					{passwordMessage.text}
				</div>
			{/if}
			<div>
				<label for="old_password" class="mb-2 block text-sm font-medium text-gray-700">Password Lama</label>
				<input type="password" id="old_password" bind:value={oldPassword} class="w-full rounded-md border border-gray-300 p-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" required />
			</div>
			<div>
				<label for="new_password" class="mb-2 block text-sm font-medium text-gray-700">Password Baru</label>
				<input type="password" id="new_password" bind:value={newPassword} class="w-full rounded-md border border-gray-300 p-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" required />
			</div>
			<button type="submit" disabled={isPasswordLoading} class="rounded-md bg-indigo-600 px-4 py-2 text-white shadow-sm hover:bg-indigo-700 disabled:cursor-not-allowed disabled:bg-indigo-400">
				{isPasswordLoading ? 'Menyimpan...' : 'Ubah Password'}
			</button>
		</form>
	</div>
</div>
