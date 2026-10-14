<script lang="ts">
	import 'prosekit/basic/style.css';
	import 'prosekit/basic/typography.css';
	import { defineBasicExtension } from 'prosekit/basic';
	import { createEditor } from 'prosekit/core';
	import { ProseKit } from 'prosekit/svelte';
	import Toolbar from '#lib/Toolbar.svelte';

	interface Props {
		content: string;
		title: string;
		saving?: boolean;
		onSave?: () => void;
	}

	let { content = $bindable(), title = $bindable(), saving = false, onSave }: Props = $props();

	const extension = defineBasicExtension();

	// Always start with an empty editor — loading markdown back into ProseKit
	// requires schema-aware conversion which is not yet implemented.
	const editor = createEditor({ extension });
</script>

<ProseKit {editor}>
	<div class="editor-shell">
		<div class="editor-topbar">
			<Toolbar {editor} onDocChange={(markdown) => { content = markdown; }} />
			<div class="topbar-actions">
				<button
					class="save-btn"
					class:save-btn--saving={saving}
					disabled={saving}
					onclick={onSave}
					type="button"
				>
					<i class="ph-duotone ph-floppy-disk"></i>
					{saving ? 'Saving…' : 'Save'}
				</button>
			</div>
		</div>
		<div class="editor-scroll">
			<input
				class="title-input"
				bind:value={title}
				placeholder="Untitled"
				type="text"
			/>
			<div {@attach editor.mount} class="editor"></div>
		</div>
	</div>
</ProseKit>

<style>
	.editor-shell {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	/* ── Top bar: toolbar + save button ── */
	.editor-topbar {
		display: flex;
		align-items: center;
		border-bottom: 1px solid #E5E5E5;
		flex-shrink: 0;
	}

	.topbar-actions {
		display: flex;
		align-items: center;
		padding: 0 12px;
		border-left: 1px solid #E5E5E5;
	}

	.save-btn {
		display: flex;
		align-items: center;
		gap: 5px;
		padding: 5px 12px;
		border: 1px solid #0F6CBD;
		border-radius: 6px;
		background: #0F6CBD;
		color: #fff;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		font-size: 12px;
		font-weight: 600;
		cursor: pointer;
		white-space: nowrap;
		transition: background 120ms ease, opacity 120ms ease;
	}

	.save-btn:hover:not(:disabled) { background: #0D5EA8; }
	.save-btn:active:not(:disabled) { transform: scale(0.97); transition-duration: 80ms; }
	.save-btn:focus-visible { box-shadow: 0 0 0 2px #005FB8; }
	.save-btn--saving { opacity: 0.65; cursor: not-allowed; }
	.save-btn:disabled { cursor: not-allowed; }

	/* ── Scrollable area ── */
	.editor-scroll {
		flex: 1;
		overflow-y: auto;
		padding: 40px 40px 48px;
		/* Must be a positioned block so the ProseMirror min-height:100% resolves */
		display: flex;
		flex-direction: column;
	}

	/* ── Title input ── */
	.title-input {
		display: block;
		width: 100%;
		max-width: 680px;
		margin: 0 auto 8px;
		padding: 4px 2px;
		border: none;
		border-bottom: 2px solid transparent;
		/* Undo @tailwindcss/forms global input overrides */
		border-radius: 0;
		box-shadow: none;
		background: transparent;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		font-size: 28px;
		font-weight: 700;
		color: #1A1A1A;
		letter-spacing: -0.02em;
		outline: none;
		transition: border-color 150ms ease;
	}

	.title-input:focus {
		border-bottom-color: #0F6CBD;
		/* Undo @tailwindcss/forms focus ring */
		box-shadow: none;
		outline: none;
	}
	.title-input::placeholder { color: #AAAAAA; }

	/* ── Editor ── */
	.editor {
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		color: #1A1A1A;
		/* Grow to fill the remaining scroll area so ProseMirror min-height:100% works */
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.editor :global(.ProseMirror) {
		max-width: 680px;
		width: 100%;
		margin: 0 auto;
		flex: 1;
		/* min-height fills the full visible scroll area so clicking anywhere focuses */
		min-height: 100%;
		outline: none;
		--prosekit-hue: 210;
	}

	.editor :global(.ProseMirror h1),
	.editor :global(.ProseMirror h2),
	.editor :global(.ProseMirror h3),
	.editor :global(.ProseMirror h4),
	.editor :global(.ProseMirror h5),
	.editor :global(.ProseMirror h6) {
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		color: #1A1A1A;
		letter-spacing: -0.01em;
	}

	.editor :global(.ProseMirror p) {
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		color: #1A1A1A;
		font-size: 16px;
		line-height: 1.65;
	}

	.editor :global(.ProseMirror code) {
		background: #F0F4F9;
		color: #0F6CBD;
		border-radius: 4px;
		padding: 1px 5px;
		font-size: 0.88em;
	}

	.editor :global(.ProseMirror blockquote) { color: #606060; border-left-color: #0F6CBD; }
	.editor :global(.ProseMirror a) { color: #0F6CBD; }
	.editor :global(.ProseMirror hr) { border-color: #E5E5E5; }
</style>

