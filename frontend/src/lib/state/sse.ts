import { writable} from 'svelte/store'
export const cachedEmitters: Map<string, (eventName: string, data: string) => import('sveltekit-sse').Unsafe<void>> = new Map();
export const cachedKeys: Map<string, boolean> = new Map();
export const clientKey = writable("")
export const messages = writable("")