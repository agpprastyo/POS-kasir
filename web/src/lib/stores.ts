import { writable } from 'svelte/store';
import type { Profile } from '$lib/types';


export const userProfile = writable<Profile | any >(null);

