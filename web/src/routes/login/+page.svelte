<script lang="ts">
	import { goto } from '$app/navigation';
	import { userProfile } from '$lib/stores';
	import '../../app.css'


	let email = '';
	let password = '';


	let isLoading = false;
	let errorMessage = '';

	async function handleLogin() {
		isLoading = true;
		errorMessage = '';

		try {
			const response = await fetch('http://127.0.0.1:8000/api/v1/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				credentials: 'include',
				body: JSON.stringify({
					email,
					password
				})
			});

			const result = await response.json();

			if (!response.ok) {
				throw new Error(result.message || 'Terjadi kesalahan saat login.');
			}

			userProfile.set(result.data);
			await goto('/');

		} catch (error: any) {
			errorMessage = error.message;
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<div class="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
		<h2 class="mb-6 text-center text-3xl font-bold text-gray-800">Selamat Datang!</h2>
		<p class="mb-6 text-center text-gray-500">Silakan masuk untuk melanjutkan</p>

		<form on:submit|preventDefault={handleLogin}>
			{#if errorMessage}
				<div class="mb-4 rounded-md bg-red-100 p-3 text-center text-sm text-red-700">
					{errorMessage}
				</div>
			{/if}

			<div class="mb-4">
				<label for="email" class="mb-2 block text-sm font-medium text-gray-700">Email</label>
				<input
					type="email"
					id="email"
					class="w-full rounded-md border border-gray-300 p-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
					placeholder="anda@email.com"
					bind:value={email}
					required
					disabled={isLoading}
				/>
			</div>

			<div class="mb-6">
				<label for="password" class="mb-2 block text-sm font-medium text-gray-700">Password</label>
				<input
					type="password"
					id="password"
					class="w-full rounded-md border border-gray-300 p-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
					placeholder="••••••••"
					bind:value={password}
					required
					disabled={isLoading}
				/>
			</div>

			<div>
				<button
					type="submit"
					class="flex w-full justify-center rounded-md bg-indigo-600 px-4 py-3 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-700 disabled:cursor-not-allowed disabled:bg-indigo-400"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="h-5 w-5 animate-spin text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						<span class="ml-2">Memproses...</span>
					{:else}
						Masuk
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>