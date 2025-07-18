<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	// Impor store dari file terpusat
	import { userProfile } from '$lib/stores';

	import '../../app.css'

	// --- **FIX:** Definisikan base URL API di satu tempat ---
	const API_BASE_URL = 'http://127.0.0.1:8000';

	// --- Definisi Menu Navigasi Berdasarkan Peran ---
	const menuItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: 'M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.25 4.5h13.5c.828 0 1.5.672 1.5 1.5v10.5c0 .828-.672 1.5-1.5 1.5H5.25c-.828 0-1.5-.672-1.5-1.5V6c0-.828.672-1.5 1.5-1.5z', roles: ['admin', 'manager'] },
		{ href: '/kasir', label: 'Kasir', icon: 'M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c.51 0 .962-.343 1.087-.835l.383-1.437M7.5 14.25L5.106 5.106A2.25 2.25 0 002.869 3H2.25', roles: ['admin', 'manager', 'cashier'] },
		{ href: '/pesanan', label: 'Pesanan', icon: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z', roles: ['admin', 'manager', 'cashier'] },
		{ href: '/produk', label: 'Produk', icon: 'M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z', roles: ['admin', 'manager'] },
		{ href: '/pengguna', label: 'Pengguna', icon: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.003c0 1.113.285 2.16.786 3.07M15 19.128c-1.33.621-2.839.621-4.277 0M6.75 16.128c-1.33.621-2.839.621-4.277 0a4.125 4.125 0 00-7.533 2.493 9.337 9.337 0 004.121.952A9.38 9.38 0 006.75 19.128v-.003c0-1.113.285-2.16.786-3.07M6.75 16.128v.003c0 1.113-.285-2.16-.786-3.07m0 0a4.125 4.125 0 00-7.533-2.493', roles: ['admin'] },
	];

	let isProfileMenuOpen = false;

	// --- Lifecycle & Data Fetching ---
	onMount(async () => {
		try {
			// Cek sesi dengan mengambil data profil.
			const response = await fetch(`${API_BASE_URL}/api/v1/auth/me`, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
				credentials: 'include'
			});

			if (!response.ok) {
				console.error(response);
				throw new Error('Sesi tidak valid');
			}

			const result = await response.json();
			userProfile.set(result.data);

		} catch (error) {
			console.error("Authentication check failed:", error);
			await goto('/login');
		}
	});

	// --- Functions ---
	async function handleLogout() {
		try {
			await fetch(`${API_BASE_URL}/api/v1/auth/logout`, {
				method: 'POST',
				credentials: 'include'
			});
		} catch (error) {
			console.error("Logout failed:", error);
		} finally {
			userProfile.set(null);
			await goto('/login');
		}
	}
</script>

{#if $userProfile}
	<div class="flex h-screen bg-gray-100 font-sans">
		<!-- Sidebar Navigasi -->
		<aside class="w-64 flex-shrink-0 bg-gray-800 text-white">
			<div class="flex h-16 items-center justify-center px-4 text-2xl font-bold">
				UMKM POS
			</div>
			<nav class="mt-4">
				{#each menuItems as item}
					<!-- Tampilkan menu hanya jika peran pengguna diizinkan -->
					{#if item.roles.includes($userProfile.role)}
						<a
							href={item.href}
							class="flex items-center px-6 py-3 transition-colors duration-200 hover:bg-gray-700"
							class:bg-gray-900={$page.url.pathname.startsWith(item.href)}
						>
							<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-6 w-6">
								<path stroke-linecap="round" stroke-linejoin="round" d={item.icon} />
							</svg>
							<span class="ml-4">{item.label}</span>
						</a>
					{/if}
				{/each}
			</nav>
		</aside>

		<!-- Konten Utama -->
		<div class="flex flex-1 flex-col overflow-hidden">
			<!-- Header -->
			<header class="flex items-center justify-end bg-white px-6 py-3 shadow-md">
				<div class="relative">
					<button on:click={() => isProfileMenuOpen = !isProfileMenuOpen} class="flex items-center space-x-2">
						<img class="h-10 w-10 rounded-full object-cover" src={$userProfile.avatar || `https://ui-avatars.com/api/?name=${$userProfile.username}&background=random`} alt="Avatar" />
						<span class="font-medium text-gray-700">{$userProfile.username}</span>
					</button>

					<!-- Dropdown Menu Profil -->
					{#if isProfileMenuOpen}
						<div class="absolute right-0 mt-2 w-48 rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5">
							<a href="/pengaturan" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Pengaturan</a>
							<button on:click={handleLogout} class="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100">
								Logout
							</button>
						</div>
					{/if}
				</div>
			</header>

			<!-- Area Konten Halaman -->
			<main class="flex-1 overflow-y-auto p-6">
				<slot />
			</main>
		</div>
	</div>
{:else}
	<!-- Tampilkan loading state atau skeleton screen saat data profil sedang diambil -->
	<div class="flex h-screen items-center justify-center bg-gray-100">
		<div class="text-center">
			<svg class="mx-auto h-12 w-12 animate-spin text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			<p class="mt-4 text-gray-600">Memuat sesi...</p>
		</div>
	</div>
{/if}
