import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

import type { CookieSerializeOptions } from 'cookie';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export const load: PageServerLoad = async ({ cookies }) => {
	const token = cookies.get('access_token');
	if (token) {
		throw redirect(303, '/');
	}
	return {};
};

export const actions: Actions = {
	default: async ({ cookies, request }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');

		if (!email || !password) {
			return fail(400, {
				error: 'Email dan password tidak boleh kosong.',
			});
		}

		try {
			const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password }),
			});

			if (!response.ok) {
				const errorResult = await response.json();
				return fail(response.status, {
					email: email.toString(),
					error: errorResult.message || 'Kredensial tidak valid.',
				});
			}

			const setCookieHeader = response.headers.get('set-cookie');
			if (setCookieHeader) {

				const [cookiePair, ...directives] = setCookieHeader.split(';').map(d => d.trim());
				const [name, value] = cookiePair.split('=');

				const options: CookieSerializeOptions = {
					path: '/',
				};

				for (const directive of directives) {
					const [key, val] = directive.split('=').map(d => d.trim());
					const lowerKey = key.toLowerCase();

					if (lowerKey === 'expires') {
						options.expires = new Date(val);
					} else if (lowerKey === 'path') {
						options.path = val;
					} else if (lowerKey === 'samesite') {
						options.sameSite = val.toLowerCase() as 'lax' | 'strict' | 'none';
					} else if (lowerKey === 'max-age') {
						options.maxAge = parseInt(val, 10);
					} else if (lowerKey === 'secure') {
						options.secure = true;
					} else if (lowerKey === 'httponly') {
						options.httpOnly = true;
					}
				}


			cookies.set(name, value, { ...options, path: options.path ?? '/' });
			}

		} catch (error) {
			console.error('Error saat login:', error);
			return fail(500, {
				error: 'Terjadi masalah pada server. Silakan coba lagi.'
			});
		}

		throw redirect(303, '/');
	},
};
