// src/lib/stores.ts
import { writable } from 'svelte/store';

// Definisikan store di sini agar bisa diimpor di mana saja
export const userProfile = writable<any>(null);