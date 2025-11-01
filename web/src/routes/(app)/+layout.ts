import { redirect } from '@sveltejs/kit';
import type {  LayoutLoad } from './$types';
import { userProfile } from '$lib/stores';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: LayoutLoad = async ({ fetch: eventFetch }) => {
	try {

		const response = await eventFetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/me`, {
			method: 'GET',
			headers: { 'Content-Type': 'application/json'},
			credentials: 'include',
		});

		console.log("Checking auth session...");
		console.log("Response status:", response.status);

		if (!response.ok) {
			throw redirect(303, '/login');
		}

		const result = await response.json();
		const profile = result.data;

		userProfile.set(profile);

		return { profile };

	} catch (error) {

		if (error instanceof Response && error.status === 303) {
			throw error;
		}
		console.error("Layout auth check failed:", error);
		throw redirect(303, '/login');
	}
};
