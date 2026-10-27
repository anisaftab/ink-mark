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
		onContentChange?: (value: string) => void;
		onTitleChange?: (value: string) => void;
	}

	let { content, title, saving = false, onSave, onContentChange, onTitleChange }: Props = $props();

	const extension = defineBasicExtension();

	// Always start with an empty editor — loading markdown back into ProseKit
	// requires schema-aware conversion which is not yet implemented.
	const editor = createEditor({ extension });
</script>

<ProseKit {editor}>
	<div class="editor-shell">
		<div class="editor-topbar">
			<Toolbar {editor} onDocChange={(markdown) => { onContentChange?.(markdown); }} />
		</div>
		<div class="editor-scroll">
			<input
				class="title-input"
				value={title}
				oninput={(e) => onTitleChange?.((e.currentTarget as HTMLInputElement).value)}
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

	/* ── Top bar: toolbar only ── */
	.editor-topbar {
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--surface-panel, #fdfcfa);
		border-bottom: 1px solid var(--border, rgba(28,25,23,0.08));
		flex-shrink: 0;
	}

	/* ── Scrollable area ── */
	.editor-scroll {
		flex: 1;
		overflow-y: auto;
		padding: 40px 40px 48px;
		/* Must be a positioned block so the ProseMirror min-height:100% resolves */
		display: flex;
		flex-direction: column;
	}

	.title-input {
		display: block;
		width: 100%;
		max-width: 680px;
		margin: 0 auto 8px;
		padding: 4px 2px;
		border: none;
		border-bottom: 2px solid transparent;
		border-radius: 0;
		box-shadow: none;
		background: transparent;
		font-family: var(--font, 'Segoe UI Variable', system-ui, sans-serif);
		font-size: 28px;
		font-weight: 700;
		color: var(--ink, #1c1917);
		letter-spacing: -0.02em;
		outline: none;
		transition: border-color 150ms ease;
	}

	.title-input:focus {
		border-bottom-color: var(--accent, #0F6CBD);
		box-shadow: none;
		outline: none;
	}
	.title-input::placeholder { color: var(--ink-subtle, #a8a29e); }

	/* ── Editor ── */
	.editor {
		font-family: var(--font, 'Segoe UI Variable', system-ui, sans-serif);
		color: var(--ink, #1c1917);
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.editor :global(.ProseMirror) {
		max-width: 680px;
		width: 100%;
		margin: 0 auto;
		flex: 1;
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
		font-family: var(--font, 'Segoe UI Variable', system-ui, sans-serif);
		color: var(--ink, #1c1917);
		letter-spacing: -0.01em;
	}

	.editor :global(.ProseMirror p) {
		font-family: var(--font, 'Segoe UI Variable', system-ui, sans-serif);
		color: var(--ink, #1c1917);
		font-size: 16px;
		line-height: 1.65;
	}

	.editor :global(.ProseMirror code) {
		background: rgba(15, 108, 189, 0.08);
		color: var(--accent, #0F6CBD);
		border-radius: 4px;
		padding: 1px 5px;
		font-size: 0.88em;
	}

	.editor :global(.ProseMirror blockquote) { color: var(--ink-muted, #57534e); border-left-color: var(--accent, #0F6CBD); }
	.editor :global(.ProseMirror a) { color: var(--accent, #0F6CBD); }
	.editor :global(.ProseMirror hr) { border-color: var(--border, rgba(28,25,23,0.08)); }
</style>

