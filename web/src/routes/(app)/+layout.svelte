<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { userProfile } from '$lib/stores';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import type { PageData } from './$types';
	import { clickOutside } from '$lib/clickOutside';


	export let data: PageData;

	const menuItems = [
		{ href: '/', label: 'Dashboard', icon: 'M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.25 4.5h13.5c.828 0 1.5.672 1.5 1.5v10.5c0 .828-.672 1.5-1.5 1.5H5.25c-.828 0-1.5-.672-1.5-1.5V6c0-.828.672-1.5 1.5-1.5z', roles: ['admin', 'manager'] },
		{ href: '/cashier', label: 'Kasir', icon: 'M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c.51 0 .962-.343 1.087-.835l.383-1.437M7.5 14.25L5.106 5.106A2.25 2.25 0 002.869 3H2.25', roles: ['admin', 'manager', 'cashier'] },
		{ href: '/order', label: 'Pesanan', icon: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z', roles: ['admin', 'manager', 'cashier'] },
		{ href: '/product', label: 'Produk', icon: 'M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z', roles: ['admin', 'manager'] },
		{ href: '/category', label: 'Kategori', icon: 'M12 3.75a8.25 8.25 0 100 16.5 8.25 8.25 0 000-16.5zM12 9a3 3 0 110 6 3 3 0 010-6zM12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75S21.75 17.385 21.75 12 17.385 2.25 12 2.25z', roles: ['admin', 'manager'] },
		{ href: '/user', label: 'Pengguna', icon: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.003c0 1.113.285 2.16.786 3.07M15 19.128c-1.33.621-2.839.621-4.277 0M6.75 16.128c-1.33.621-2.839.621-4.277 0a4.125 4.125 0 00-7.533 2.493 9.337 9.337 0 004.121.952A9.38 9.38 0 006.75 19.128v-.003c0-1.113.285-2.16.786-3.07M6.75 16.128v.003c0 1.113-.285-2.16-.786-3.07m0 0a4.125 4.125 0 00-7.533-2.493', roles: ['admin'] },
	];

	let isProfileMenuOpen = false;

	async function handleLogout() {
		try {
			await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/logout`, {
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


<div class="flex h-screen bg-white font-sans">
	<aside class="flex w-64 flex-col border-r border-gray-200 bg-white">
		<div class="flex h-16 shrink-0 items-center px-6">
			<span class="text-xl font-bold text-indigo-600">UMKM POS</span>
		</div>

		<!-- Scrollable nav -->
		<div class="flex grow flex-col overflow-y-auto px-6">
			<nav class="flex flex-1 flex-col">
				<ul role="list" class="flex flex-1 flex-col gap-y-7">
					<li>
						<ul role="list" class="-mx-2 space-y-1">
							{#each menuItems as item (item.href)}
								{#if data.profile && item.roles.includes(data.profile.role)}
									<li>
										<a
											href={item.href}
											class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold
											hover:bg-gray-50 hover:text-indigo-600
											{ $page.url.pathname === item.href
												? 'bg-gray-50 text-indigo-600'
												: 'text-gray-700' }"
										>
											<svg
												class="size-6 shrink-0
												{ $page.url.pathname === item.href
													? 'text-indigo-600'
													: 'text-gray-400 group-hover:text-indigo-600' }"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="1.5"
												stroke="currentColor"
												aria-hidden="true"
											>
												<path stroke-linecap="round" stroke-linejoin="round" d={item.icon} />
											</svg>
											<span class="truncate">{item.label}</span>
										</a>
									</li>
								{/if}
							{/each}
						</ul>
					</li>
					<li class="-mx-6 mt-auto">
						<div class="relative px-6 py-3">
							<button
								on:click={() => (isProfileMenuOpen = !isProfileMenuOpen)}
								class="flex w-full items-center gap-x-4 text-sm font-semibold text-gray-900 hover:bg-gray-100 transition-colors rounded-lg p-2"
							>
								<img
									class="size-8 rounded-full bg-gray-100 object-cover"
									src={data.profile.avatar || `https://ui-avatars.com/api/?name=${data.profile.username}&background=random`}
									alt="Avatar"
								/>
								<span class="truncate">{data.profile.username}</span>
							</button>

							{#if isProfileMenuOpen}
								<div
									use:clickOutside={() => (isProfileMenuOpen = false)}
									class="absolute bottom-16 right-6 z-50 w-48 rounded-md bg-white py-1 shadow ring-1 ring-black ring-opacity-5 transition-all"
								>
									<a
										href="/settings"
										on:click={() => (isProfileMenuOpen = !isProfileMenuOpen)}
										class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
									>
										Pengaturan
									</a>
									<button
										on:click={handleLogout}
										class="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 transition-colors"
									>
										Logout
									</button>
								</div>
							{/if}
						</div>
					</li>
				</ul>
			</nav>
		</div>


	</aside>



	<div class="flex flex-1 flex-col overflow-hidden">
		<main class="flex-1 overflow-y-auto p-6">
			<slot />
		</main>
	</div>
</div>
