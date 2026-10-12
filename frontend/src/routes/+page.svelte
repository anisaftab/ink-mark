<script lang="ts">
	import Editor from '#lib/Editor.svelte';
	import { NewDocument, SaveDocument } from '#lib/wailsjs/go/main/App';
	import type { Document } from '#lib/types';
	import { newDocumentTrigger } from '#lib/stores';

	let activeDocument = $state<Document | null>(null);
	let saving = $state(false);
	let error = $state<string | null>(null);

	// React to the sidebar New Document button.
	$effect(() => {
		const unsub = newDocumentTrigger.subscribe((n: number) => {
			if (n > 0) handleNewDocument();
		});
		return unsub;
	});

	async function handleNewDocument(): Promise<void> {
		error = null;
		try {
			const doc = await NewDocument();
			// User may have cancelled the save dialog — doc will have empty id.
			if (doc.id) activeDocument = doc;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create document';
		}
	}

	async function handleSave(): Promise<void> {
		if (!activeDocument) return;
		saving = true;
		error = null;
		try {
			await SaveDocument(
				activeDocument.id,
				activeDocument.title,
				activeDocument.content
			);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save document';
		} finally {
			saving = false;
		}
	}
</script>

{#if activeDocument}
	<Editor
		bind:content={activeDocument.content}
		bind:title={activeDocument.title}
		onSave={handleSave}
		{saving}
	/>
{:else}
	<div class="empty-state">
		<i class="ph-duotone ph-notebook empty-icon"></i>
		<p class="empty-title">No document open</p>
		<button class="new-btn" onclick={handleNewDocument} type="button">
			<i class="ph-duotone ph-plus"></i>
			New Document
		</button>
		{#if error}
			<p class="error">{error}</p>
		{/if}
	</div>
{/if}

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 12px;
		user-select: none;
	}

	.empty-icon {
		font-size: 48px;
		color: #0F6CBD;
		opacity: 0.4;
	}

	.empty-title {
		font-size: 15px;
		color: #606060;
		margin: 0;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
	}

	.new-btn {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		border: 1px solid #0F6CBD;
		border-radius: 6px;
		background: rgba(15, 108, 189, 0.08);
		color: #0F6CBD;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		font-size: 13px;
		font-weight: 600;
		cursor: pointer;
		transition: background 120ms ease;
	}

	.new-btn:hover { background: rgba(15, 108, 189, 0.14); }
	.new-btn:active { background: rgba(15, 108, 189, 0.20); transform: scale(0.98); }

	.error {
		font-size: 12px;
		color: #C42B1C;
		margin: 0;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
	}
</style>
