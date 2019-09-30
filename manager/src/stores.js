import {writable} from 'svelte/store';

export const token = writable(localStorage.getItem('token'));
export const authenticated = writable(false);
export const authenticating = writable(false);
export const profile = writable(null);
