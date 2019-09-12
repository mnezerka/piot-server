import {writable} from 'svelte/store';

export const token = writable(localStorage.getItem('token'));
export const authenticated = writable(false);
export const profile = writable(null);
//export const test_plans = writable(null);
//export const test_plan = writable(null);
