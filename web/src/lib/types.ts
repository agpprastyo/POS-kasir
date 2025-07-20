export type ErrorResponse = {
	message: string;
	errors?: Record<string, string[]>;
}
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
