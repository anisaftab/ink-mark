<script lang="ts">
	import type { BasicExtension } from 'prosekit/basic';
	import { useEditor, useEditorDerivedValue } from 'prosekit/svelte';

	const editor = useEditor<BasicExtension>({ update: true });
	const state = useEditorDerivedValue<BasicExtension, Record<string, { active: boolean; canExec: boolean }>>((e) => ({
		bold:        { active: e.marks.bold?.isActive() ?? false,                    canExec: e.commands.toggleBold?.canExec() ?? false },
		italic:      { active: e.marks.italic?.isActive() ?? false,                  canExec: e.commands.toggleItalic?.canExec() ?? false },
		underline:   { active: e.marks.underline?.isActive() ?? false,               canExec: e.commands.toggleUnderline?.canExec() ?? false },
		strike:      { active: e.marks.strike?.isActive() ?? false,                  canExec: e.commands.toggleStrike?.canExec() ?? false },
		code:        { active: e.marks.code?.isActive() ?? false,                    canExec: e.commands.toggleCode?.canExec() ?? false },
		h1:          { active: e.nodes.heading?.isActive({ level: 1 }) ?? false,     canExec: e.commands.toggleHeading?.canExec({ level: 1 }) ?? false },
		h2:          { active: e.nodes.heading?.isActive({ level: 2 }) ?? false,     canExec: e.commands.toggleHeading?.canExec({ level: 2 }) ?? false },
		h3:          { active: e.nodes.heading?.isActive({ level: 3 }) ?? false,     canExec: e.commands.toggleHeading?.canExec({ level: 3 }) ?? false },
		bulletList:  { active: e.nodes.list?.isActive({ kind: 'bullet' }) ?? false,  canExec: e.commands.toggleList?.canExec({ kind: 'bullet' }) ?? false },
		orderedList: { active: e.nodes.list?.isActive({ kind: 'ordered' }) ?? false, canExec: e.commands.toggleList?.canExec({ kind: 'ordered' }) ?? false },
		blockquote:  { active: e.nodes.blockquote?.isActive() ?? false,              canExec: e.commands.toggleBlockquote?.canExec() ?? false },
	}));
</script>

<div class="toolbar" role="toolbar" aria-label="Text formatting">
	<div class="toolbar-group">
		<button class="btn" class:active={$state.h1.active} disabled={!$state.h1.canExec} onclick={() => $editor.commands.toggleHeading({ level: 1 })} title="Heading 1" type="button">H1</button>
		<button class="btn" class:active={$state.h2.active} disabled={!$state.h2.canExec} onclick={() => $editor.commands.toggleHeading({ level: 2 })} title="Heading 2" type="button">H2</button>
		<button class="btn" class:active={$state.h3.active} disabled={!$state.h3.canExec} onclick={() => $editor.commands.toggleHeading({ level: 3 })} title="Heading 3" type="button">H3</button>
	</div>
	<div class="sep"></div>
	<div class="toolbar-group">
		<button class="btn icon" class:active={$state.bold.active}      disabled={!$state.bold.canExec}      onclick={() => $editor.commands.toggleBold()}            title="Bold (Ctrl+B)"  type="button"><i class="ph-duotone ph-text-b"></i></button>
		<button class="btn icon" class:active={$state.italic.active}    disabled={!$state.italic.canExec}    onclick={() => $editor.commands.toggleItalic()}          title="Italic (Ctrl+I)" type="button"><i class="ph-duotone ph-text-italic"></i></button>
		<button class="btn icon" class:active={$state.underline.active} disabled={!$state.underline.canExec} onclick={() => $editor.commands.toggleUnderline()}       title="Underline"      type="button"><i class="ph-duotone ph-text-underline"></i></button>
		<button class="btn icon" class:active={$state.strike.active}    disabled={!$state.strike.canExec}    onclick={() => $editor.commands.toggleStrike()}          title="Strikethrough"  type="button"><i class="ph-duotone ph-text-strikethrough"></i></button>
		<button class="btn icon" class:active={$state.code.active}      disabled={!$state.code.canExec}      onclick={() => $editor.commands.toggleCode()}            title="Inline code"    type="button"><i class="ph-duotone ph-code"></i></button>
	</div>
	<div class="sep"></div>
	<div class="toolbar-group">
		<button class="btn icon" class:active={$state.bulletList.active}  disabled={!$state.bulletList.canExec}  onclick={() => $editor.commands.toggleList({ kind: 'bullet' })}  title="Bullet list"   type="button"><i class="ph-duotone ph-list-bullets"></i></button>
		<button class="btn icon" class:active={$state.orderedList.active} disabled={!$state.orderedList.canExec} onclick={() => $editor.commands.toggleList({ kind: 'ordered' })} title="Numbered list" type="button"><i class="ph-duotone ph-list-numbers"></i></button>
		<button class="btn icon" class:active={$state.blockquote.active}  disabled={!$state.blockquote.canExec}  onclick={() => $editor.commands.toggleBlockquote()}              title="Blockquote"    type="button"><i class="ph-duotone ph-quotes"></i></button>
	</div>
</div>

<style>
	.toolbar {
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 8px 16px;
		background: #F8F8F8;
		border-bottom: 1px solid #E5E5E5;
		flex-shrink: 0;
	}

	.toolbar-group { display: flex; align-items: center; gap: 2px; }

	.sep {
		width: 1px;
		height: 18px;
		background: #D0D0D0;
		margin: 0 8px;
		flex-shrink: 0;
	}

	.btn {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 30px;
		height: 28px;
		padding: 0 8px;
		border: 1px solid transparent;
		border-radius: 6px;
		background: transparent;
		color: #606060;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		font-size: 12px;
		font-weight: 600;
		letter-spacing: -0.02em;
		cursor: pointer;
		outline: none;
		transition: background 120ms ease, color 120ms ease, border-color 120ms ease;
	}

	.icon { font-size: 16px; padding: 0; }

	.btn:hover:not(:disabled) {
		background: #F0F0F0;
		color: #1A1A1A;
		border-color: #E5E5E5;
	}

	.btn:active:not(:disabled) {
		background: #E7E7E7;
		transform: scale(0.97);
		transition-duration: 80ms;
	}

	.btn:focus-visible { box-shadow: 0 0 0 2px #005FB8; }

	.btn.active {
		background: rgba(15, 108, 189, 0.10);
		color: #0F6CBD;
		border-color: rgba(15, 108, 189, 0.25);
	}

	.btn.active:hover:not(:disabled) { background: rgba(15, 108, 189, 0.16); }

	.btn:disabled { opacity: 0.38; cursor: not-allowed; }
</style>

