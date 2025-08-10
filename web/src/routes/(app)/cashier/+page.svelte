<script lang="ts">
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';
	import type { Product, ProductOption, OrderItemRequest } from '$lib/types';
	import { createOrder } from '$lib/api/order';

	export let data: PageData;

	// --- State untuk UI ---
	let search = data.queryParams?.search || '';
	let category_id = data.queryParams?.category_id || '';

	// --- State Keranjang Belanja ---
	let cart: (OrderItemRequest & { name: string; price: number; selectedOptionsData: ProductOption[] })[] = [];
	let cartTotal = 0;

	// --- State Modal Opsi ---
	let showOptionsModal = false;
	let selectedProduct: Product | null = null;
	let selectedOptions: { [key: string]: boolean } = {};

	// --- Fungsi ---
	function formatCurrency(value: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value);
	}

	function handleFilterSubmit() {
		const params = new URLSearchParams();
		if (search) params.set('search', search);
		if (category_id) params.set('category_id', category_id);
		goto(`?${params.toString()}`, { keepFocus: true, noScroll: true, replaceState: true });
	}

	// Fungsi baru untuk menangani klik pada tab kategori
	function selectCategory(id: string | number) {
		category_id = String(id);
		handleFilterSubmit();
	}

	function openOptionsModal(product: Product) {
		selectedProduct = product;
		selectedOptions = {}; // Reset pilihan
		showOptionsModal = true;
	}

	function addToCart(product: Product, chosenOptions: ProductOption[] = []) {
		const existingItemIndex = cart.findIndex(item =>
			item.product_id === product.id &&
			JSON.stringify(item.options.map(o => o.product_option_id).sort()) === JSON.stringify(chosenOptions.map(o => o.id).sort())
		);

		let itemPrice = product.price;
		chosenOptions.forEach(opt => itemPrice += opt.additional_price);

		if (existingItemIndex > -1) {
			cart[existingItemIndex].quantity++;
		} else {
			cart = [...cart, {
				product_id: product.id,
				quantity: 1,
				options: chosenOptions.map(o => ({ product_option_id: o.id! })),
				name: product.name,
				price: itemPrice,
				selectedOptionsData: chosenOptions
			}];
		}
		updateCartTotal();
	}

	function handleAddToCartFromModal() {
		if (!selectedProduct) return;
		const chosenOptions = selectedProduct.options?.filter(opt => selectedOptions[opt.id!]) || [];
		addToCart(selectedProduct, chosenOptions);
		showOptionsModal = false;
	}

	function handleProductClick(product: Product) {
		if (product.options && product.options.length > 0) {
			openOptionsModal(product);
		} else {
			addToCart(product);
		}
	}

	function updateCartTotal() {
		cartTotal = cart.reduce((total, item) => total + (item.price * item.quantity), 0);
	}

	function updateQuantity(index: number, change: number) {
		const newQuantity = cart[index].quantity + change;
		if (newQuantity > 0) {
			cart[index].quantity = newQuantity;
		} else {
			cart = cart.filter((_, i) => i !== index);
		}
		cart = [...cart]; // Trigger reactivity
		updateCartTotal();
	}

	async function finalizeOrder() {
		if (cart.length === 0) return;

		const orderRequest: import('$lib/types').CreateOrderRequest = {
			type: 'takeaway',
			items: cart.map(item => ({
				product_id: item.product_id,
				quantity: item.quantity,
				options: item.options
			}))
		};

		try {
			const result = await createOrder(orderRequest);
			alert(`Pesanan ${result.data.id} berhasil dibuat!`);
			// Arahkan ke halaman detail pesanan atau reset keranjang
			cart = [];
			updateCartTotal();
		} catch (error: any) {
			alert(`Gagal membuat pesanan: ${error.message}`);
		}
	}

</script>

<div class="flex h-screen ">
	<!-- Kolom Kiri: Produk -->
	<div class="flex-grow pr-6 overflow-y-auto">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-3xl font-bold text-gray-800">Kasir</h1>
			<!-- Search Bar -->
			<form on:submit|preventDefault={handleFilterSubmit} class="w-full max-w-sm">
				<div class="relative">
					<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
						<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35M17 11a6 6 0 11-12 0 6 6 0 0112 0z"></path></svg>
					</div>
					<input type="search" bind:value={search} placeholder="Cari produk..." class="block w-full rounded-lg border-gray-300 shadow-sm pl-10" />
				</div>
			</form>
		</div>

		<!-- Tab Bar Kategori -->
		<div class="mb-6 overflow-x-auto">
			<div class="flex space-x-4 ">
				<button on:click={() => selectCategory('')} class="whitespace-nowrap py-3 px-1 text-sm font-medium border-b-2 transition-colors"
								class:border-indigo-500={category_id === ''}
								class:text-indigo-600={category_id === ''}
								class:border-transparent={category_id !== ''}
								class:text-gray-500={category_id !== ''}
								class:hover:text-gray-700={category_id !== ''}
								class:hover:border-gray-300={category_id !== ''}>
					Semua Kategori
				</button>
				{#each data.categories as category (category.id)}
					<button on:click={() => selectCategory(category.id)} class="whitespace-nowrap py-3 px-1 text-sm font-medium border-b-2 transition-colors"
									class:border-indigo-500={category_id == category.id}
									class:text-indigo-600={category_id == category.id}
									class:border-transparent={category_id != category.id}
									class:text-gray-500={category_id != category.id}
									class:hover:text-gray-700={category_id != category.id}
									class:hover:border-gray-300={category_id != category.id}>
						{category.name}
					</button>
				{/each}
			</div>
		</div>

		<!-- Grid Produk -->
		<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-4 gap-4">
			{#each data.products as product (product.id)}
				<button on:click={() => handleProductClick(product)} class="text-left rounded-lg bg-white shadow-sm overflow-hidden transition-transform hover:-translate-y-1 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
					<img src={product.image_url || `https://placehold.co/300x300/e2e8f0/64748b?text=${product.name}`} alt={product.name} class=" object-cover aspect-square" />
					<div class="p-3">
						<p class="font-semibold text-gray-800 truncate">{product.name}</p>
						<p class="text-sm text-gray-600">{formatCurrency(product.price)}</p>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Kolom Kanan: Keranjang -->
	<div class="w-96 flex-shrink-0 bg-white shadow-lg p-6 flex flex-col">
		<h2 class="text-2xl font-bold border-b pb-4 mb-4">Pesanan</h2>
		<div class="flex-grow overflow-y-auto -mr-6 pr-6">
			{#if cart.length === 0}
				<div class="flex flex-col items-center justify-center h-full text-center text-gray-500">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
					</svg>
					<p class="mt-4">Keranjang masih kosong.</p>
				</div>
			{:else}
				{#each cart as item, i (i)}
					<div class="mb-4">
						<p class="font-semibold">{item.name}</p>
						{#if item.selectedOptionsData.length > 0}
							<p class="text-xs text-gray-500 pl-2">
								+ {item.selectedOptionsData.map(o => o.name).join(', ')}
							</p>
						{/if}
						<div class="flex items-center justify-between mt-1">
							<div class="flex items-center border rounded-md">
								<button on:click={() => updateQuantity(i, -1)} class="px-2 py-1 text-gray-600 hover:bg-gray-100">-</button>
								<span class="px-3 font-medium">{item.quantity}</span>
								<button on:click={() => updateQuantity(i, 1)} class="px-2 py-1 text-gray-600 hover:bg-gray-100">+</button>
							</div>
							<p class="text-sm font-semibold">{formatCurrency(item.price * item.quantity)}</p>
						</div>
					</div>
				{/each}
			{/if}
		</div>
		<div class="border-t pt-4">
			<div class="flex justify-between font-bold text-xl mb-4">
				<span>Total</span>
				<span>{formatCurrency(cartTotal)}</span>
			</div>
			<button on:click={finalizeOrder} disabled={cart.length === 0} class="w-full rounded-lg bg-green-600 py-3 text-white font-bold hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors">
				Bayar
			</button>
		</div>
	</div>
</div>

<!-- Modal Opsi Produk -->
{#if showOptionsModal && selectedProduct}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60" on:click|self={() => showOptionsModal = false}>
		<div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
			<h3 class="text-xl font-bold mb-4">Pilih Opsi untuk {selectedProduct.name}</h3>
			<div class="space-y-3 mb-6 max-h-64 overflow-y-auto">
				{#each selectedProduct.options || [] as option (option.id)}
					<label class="flex items-center justify-between p-3 rounded-lg border cursor-pointer has-[:checked]:bg-indigo-50 has-[:checked]:border-indigo-400 transition-colors">
						<div>
							<p class="font-medium">{option.name}</p>
							<p class="text-sm text-gray-600">+ {formatCurrency(option.additional_price)}</p>
						</div>
						<input type="checkbox" bind:checked={selectedOptions[option.id!]} class="h-5 w-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
					</label>
				{/each}
			</div>
			<div class="flex justify-end gap-4">
				<button on:click={() => showOptionsModal = false} class="rounded-md bg-gray-200 px-4 py-2 hover:bg-gray-300">Batal</button>
				<button on:click={handleAddToCartFromModal} class="rounded-md bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700">Tambah ke Keranjang</button>
			</div>
		</div>
	</div>
{/if}
