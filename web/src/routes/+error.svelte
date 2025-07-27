<script lang="ts">
	import { page } from '$app/stores';

	// Objek untuk memetakan status code ke pesan yang lebih ramah pengguna
	const errorDetails: { [key: number]: { title: string; message: string } } = {
		401: {
			title: 'Tidak Terotentikasi',
			message: 'Anda harus login terlebih dahulu untuk dapat mengakses halaman ini.'
		},
		403: {
			title: 'Akses Ditolak',
			message: 'Anda tidak memiliki izin yang cukup untuk melihat halaman ini.'
		},
		404: {
			title: 'Halaman Tidak Ditemukan',
			message: 'Maaf, halaman yang Anda cari tidak ada atau mungkin telah dipindahkan.'
		},
		500: {
			title: 'Kesalahan Server Internal',
			message: 'Terjadi masalah pada server kami. Tim kami telah diberitahu dan sedang bekerja untuk memperbaikinya.'
		}
	};

	// Dapatkan detail error berdasarkan status code dari $page store,
	// atau gunakan pesan default jika status code tidak dikenali.
	$: details = errorDetails[$page.status] || {
		title: 'Terjadi Kesalahan',
		message: $page.error?.message || 'Sesuatu yang tidak terduga telah terjadi. Silakan coba lagi nanti.'
	};
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 text-center">
	<div class="w-full max-w-lg flex flex-col justify-center items-center">
		<h1 class="text-8xl font-extrabold tracking-tight text-indigo-600 sm:text-9xl">
			{$page.status}
		</h1>
		<h2 class="mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
			{details.title}
		</h2>
		<p class="mt-4  text-base text-gray-500 ">
			{details.message}
		</p>
		<div class="mt-20 ">
			<a href="/" class="rounded-lg bg-indigo-600 px-6 py-3 font-semibold text-white shadow-sm transition hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
				Kembali ke Halaman Utama
			</a>
		</div>
	</div>
</div>
