<script lang="ts">
	import 'prosekit/basic/style.css';
	import 'prosekit/basic/typography.css';
	import { defineBasicExtension } from 'prosekit/basic';
	import { createEditor } from 'prosekit/core';
	import { ProseKit } from 'prosekit/svelte';
	import Toolbar from '#lib/Toolbar.svelte';

	const extension = defineBasicExtension();
	const editor = createEditor({ extension });
</script>

<ProseKit {editor}>
	<div class="editor-shell">
		<Toolbar />
		<div class="editor-scroll">
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

	.editor-scroll {
		flex: 1;
		overflow-y: auto;
		padding: 48px 40px;
	}

	.editor {
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		color: #1A1A1A;
	}

	.editor :global(.ProseMirror) {
		max-width: 680px;
		margin: 0 auto;
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

	.editor :global(.ProseMirror blockquote) {
		color: #606060;
		border-left-color: #0F6CBD;
	}

	.editor :global(.ProseMirror a) { color: #0F6CBD; }
	.editor :global(.ProseMirror hr) { border-color: #E5E5E5; }
</style>

