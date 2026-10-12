import { writable } from 'svelte/store';

/** Increment to trigger a new document creation from anywhere in the tree. */
export const newDocumentTrigger = writable(0);
