import { writable, derived } from 'svelte/store';
import type { Document } from '#lib/types';

/** Increment to trigger a new document creation from anywhere in the tree. */
export const newDocumentTrigger = writable(0);

/** All documents currently open in tabs, in order. */
export const openDocuments = writable<Document[]>([]);

/** The ID of the tab that is currently visible in the editor. */
export const activeDocumentId = writable<string | null>(null);

/** The full Document object for the active tab, or null if none is open. */
export const activeDocument = derived(
	[openDocuments, activeDocumentId],
	([$docs, $id]) => $docs.find((d) => d.id === $id) ?? null
);

/** Open a document in a tab. If it is already open, just activate it. */
export function openInTab(doc: Document): void {
	openDocuments.update((docs) => {
		const exists = docs.some((d) => d.id === doc.id);
		return exists ? docs : [...docs, doc];
	});
	activeDocumentId.set(doc.id);
}

/** Close a tab by document ID. Activates the nearest remaining tab. */
export function closeTab(id: string): void {
	openDocuments.update((docs) => {
		const idx = docs.findIndex((d) => d.id === id);
		if (idx === -1) return docs;

		const next = docs.filter((d) => d.id !== id);

		activeDocumentId.update((activeId) => {
			if (activeId !== id) return activeId;
			if (next.length === 0) return null;
			// Prefer the tab to the right, fall back to the one to the left.
			return (next[idx] ?? next[idx - 1]).id;
		});

		return next;
	});
}

/** Update the in-memory content of an open document (tracks unsaved edits). */
export function updateDocumentContent(id: string, content: string): void {
	openDocuments.update((docs) =>
		docs.map((d) => (d.id === id ? { ...d, content } : d))
	);
}

/** Update the in-memory title of an open document. */
export function updateDocumentTitle(id: string, title: string): void {
	openDocuments.update((docs) =>
		docs.map((d) => (d.id === id ? { ...d, title } : d))
	);
}

/** Flat list of all documents known to the sidebar (loaded on startup). */
export const allDocuments = writable<Document[]>([]);

