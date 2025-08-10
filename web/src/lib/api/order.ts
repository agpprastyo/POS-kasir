import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type { CreateOrderRequest, Order } from '$lib/types';

/**
 * Fungsi pembantu untuk menangani respons dari API.
 */
async function handleResponse(response: Response) {
	const result = await response.json();
	if (!response.ok) {
		throw new Error(result.message || `HTTP error! status: ${response.status}`);
	}
	return result;
}

/**
 * Membuat pesanan baru.
 * @param {CreateOrderRequest} orderData - Data pesanan baru.
 * @returns {Promise<{ data: Order }>} - Data pesanan yang baru dibuat.
 */
export async function createOrder(orderData: CreateOrderRequest): Promise<{ data: Order }> {
	const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/orders`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
		body: JSON.stringify(orderData),
	});
	return handleResponse(response);
}
