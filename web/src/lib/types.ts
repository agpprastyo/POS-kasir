
export type Profile = {
	id: string;
	username: string;
	email: string;
	created_at: string;
	updated_at: string;
	avatar: string | null;
	role: 'admin' | 'manager' | 'cashier';
	is_active: boolean;
};

export type UpdatePasswordRequest = {
	old_password: string;
	new_password: string;
};

// ... (Tipe Profile dan UpdatePasswordRequest yang sudah ada)

// Tipe untuk query parameter saat mengambil daftar pengguna
export type UserQueryParams = {
	page?: number;
	limit?: number;
	search?: string;
	role?: 'admin' | 'manager' | 'cashier' | '';
	is_active?: boolean | string;
	sortBy?: 'username' | 'email' | 'created_at';
	sortOrder?: 'asc' | 'desc';
};

// Tipe untuk informasi paginasi dari API
export type PaginationInfo = {
	current_page: number;
	total_page: number;
	total_data: number;
	per_page: number;
};

// Tipe untuk respons lengkap dari API pengguna
export type UsersApiResponse = {
	users: Profile[];
	pagination: PaginationInfo;
};


// Tipe untuk data yang dikirim saat membuat pengguna baru
export type CreateUserRequest = {
	username: string;
	email: string;
	password: string;
	role: 'admin' | 'manager' | 'cashier';
	is_active: boolean;
};

// Tipe untuk data yang dikirim saat memperbarui pengguna
export type UpdateUserRequest = {
	username: string;
	email: string;
	role: 'admin' | 'manager' | 'cashier';
};


// Tipe untuk data kategori (tetap berguna untuk filter)
export type Category = {
	id: number;
	name: string;
};

// Tipe untuk opsi produk
export type ProductOption = {
	id?: string; // Opsional saat membuat
	name: string;
	additional_price: number;
	image_url?: string; // Opsional
};

// **DIPERBARUI:** Tipe untuk data produk
export type Product = {
	id: string;
	name: string;
	category_id: number;
	category_name: string;
	price: number;
	stock: number;
	sku?: string;
	image?: string | null;
	image_url?: string | null; // Dari respons upload
	created_at?: string;
	updated_at?: string;
	options?: ProductOption[];
};

// Tipe untuk query parameter saat mengambil daftar produk
export type ProductQueryParams = {
	page?: number;
	limit?: number;
	search?: string;
	category_id?: number | string;
};

// Tipe untuk respons lengkap dari API produk
export type ProductsApiResponse = {
	products: Product[];
	pagination: PaginationInfo;
};

// Tipe untuk data yang dikirim saat membuat opsi produk
export type CreateProductOption = {
	name: string;
	additional_price: number;
};

// Tipe untuk data yang dikirim saat membuat produk baru
export type CreateProductRequest = {
	name: string;
	category_id: number;
	price: number;
	stock: number;
	options?: CreateProductOption[];
};

export type ErrorResponse = {
	message: string;
	error:   null | string;
}


// Tipe untuk data yang dikirim saat memperbarui produk utama
export type UpdateProductRequest = {
	name: string;
	category_id: number;
	price: number;
	stock: number;
};

// Tipe untuk data yang dikirim saat memperbarui opsi produk
export type UpdateProductOptionRequest = {
	name: string;
	additional_price: number;
};

// Tipe untuk data kategori dengan jumlah produk
export type CategoryWithCount = {
	id: number;
	name: string;
	product_count: number;
	created_at: string;
	updated_at: string;
};

// Tipe untuk data yang dikirim saat membuat atau memperbarui kategori
export type CategoryRequest = {
	name: string;
};

// Tipe untuk opsi yang dipilih dalam satu item pesanan
export type OrderItemOptionRequest = {
	product_option_id: string;
};

// Tipe untuk satu item dalam permintaan pembuatan pesanan
export type OrderItemRequest = {
	product_id: string;
	quantity: number;
	options: OrderItemOptionRequest[];
};

// Tipe untuk permintaan pembuatan pesanan baru
export type CreateOrderRequest = {
	type: 'dine-in' | 'takeaway';
	items: OrderItemRequest[];
};

// Tipe untuk opsi dalam item pesanan yang diterima dari API
export type OrderItemOption = {
	product_option_id: string;
	price_at_sale: number;
};

// Tipe untuk satu item dalam pesanan yang diterima dari API
export type OrderItem = {
	id: string;
	product_id: string;
	quantity: number;
	price_at_sale: number;
	subtotal: number;
	options: OrderItemOption[];
};

// Tipe untuk data pesanan lengkap yang diterima dari API
export type Order = {
	id: string;
	user_id: string;
	type: 'dine-in' | 'takeaway';
	status: 'open' | 'paid' | 'cancelled';
	gross_total: number;
	discount_amount: number;
	net_total: number;
	created_at: string;
	updated_at: string;
	items: OrderItem[];
};


