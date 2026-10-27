<script lang="ts">
	import Editor from '#lib/Editor.svelte';
	import { NewDocument, OpenDocument, SaveDocument } from '#lib/wailsjs/go/main/App';
	import {
		newDocumentTrigger,
		activeDocument,
		openInTab,
		updateDocumentContent,
		updateDocumentTitle,
		allDocuments,
		dirtyDocumentIds,
		markDirty,
		markClean,
		editorStats
	} from '#lib/stores';
	import MarkdownIt from 'markdown-it';

	const md = new MarkdownIt({ html: false, linkify: true, typographer: true });

	let saving = $state(false);
	let error = $state<string | null>(null);
	let loadingNew = $state(false);
	let loadingOpen = $state(false);
	let viewMode = $state<'edit' | 'split' | 'read'>('split');

	// Derived: rendered HTML from the active document's markdown content.
	let previewHtml = $derived(
		$activeDocument ? md.render($activeDocument.content || '') : ''
	);

	// React to the sidebar / tab-bar "New document" trigger.
	$effect(() => {
		const unsub = newDocumentTrigger.subscribe((n: number) => {
			if (n > 0) handleNewDocument();
		});
		return unsub;
	});

	// Ctrl+S → save.
	function handleKeydown(e: KeyboardEvent) {
		if ((e.ctrlKey || e.metaKey) && e.key === 's') {
			e.preventDefault();
			handleSave();
		}
	}

	async function handleNewDocument(): Promise<void> {
		if (loadingNew) return;
		loadingNew = true;
		error = null;
		try {
			const doc = await NewDocument();
			if (doc.id) {
				openInTab(doc);
				allDocuments.update((docs) => {
					const exists = docs.some((d) => d.id === doc.id);
					return exists ? docs : [doc, ...docs];
				});
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create document';
		} finally {
			loadingNew = false;
		}
	}

	async function handleOpenDocument(): Promise<void> {
		if (loadingOpen) return;
		loadingOpen = true;
		error = null;
		try {
			const doc = await OpenDocument();
			if (doc.id) {
				openInTab(doc);
				allDocuments.update((docs) => {
					const exists = docs.some((d) => d.id === doc.id);
					return exists ? docs : [doc, ...docs];
				});
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to open document';
		} finally {
			loadingOpen = false;
		}
	}

	async function handleSave(): Promise<void> {
		if (!$activeDocument) return;
		saving = true;
		error = null;
		try {
			await SaveDocument(
				$activeDocument.id,
				$activeDocument.title,
				$activeDocument.content
			);
			markClean($activeDocument.id);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save document';
		} finally {
			saving = false;
		}
	}

	function handleContentChange(id: string, value: string) {
		updateDocumentContent(id, value);
		markDirty(id);
		// Update word count from markdown content.
		const words = value.trim() ? value.trim().split(/\s+/).length : 0;
		editorStats.update((s) => ({ ...s, words }));
	}

	function handleTitleChange(id: string, value: string) {
		updateDocumentTitle(id, value);
		markDirty(id);
	}

	function handleCursorChange(line: number, col: number) {
		editorStats.update((s) => ({ ...s, line, col }));
	}
</script>


<svelte:window onkeydown={handleKeydown} />

{#if $activeDocument}
	<div class="doc-toolbar">
		<div class="breadcrumb">
			<i class="ph-duotone ph-file-text" style="font-size:15px;color:var(--ic-doc)"></i>
			<span class="breadcrumb-cur">{$activeDocument.title || 'Untitled'}</span>
		</div>
		<div class="save-status">
			{#if saving}
				<i class="ph-duotone ph-circle-notch" style="font-size:14px;color:var(--ic-settings);animation:spin 900ms linear infinite"></i>
				Saving…
			{:else if $dirtyDocumentIds.has($activeDocument.id)}
				<i class="ph-duotone ph-dot-outline" style="font-size:14px;color:var(--ic-pin)"></i>
				Unsaved
			{:else}
				<i class="ph-duotone ph-check-circle" style="font-size:14px;color:var(--ic-secure)"></i>
				Saved
			{/if}
		</div>
		<div class="toolbar-spacer"></div>
		<div class="mode-switch" role="group" aria-label="View mode">
			<button class="mode-btn" class:mode-btn--active={viewMode === 'edit'}
				onclick={() => { viewMode = 'edit'; }} type="button">
				<i class="ph-duotone ph-pencil-simple-line" style="font-size:15px;color:var(--ic-create)"></i>Edit
			</button>
			<button class="mode-btn" class:mode-btn--active={viewMode === 'split'}
				onclick={() => { viewMode = 'split'; }} type="button">
				<i class="ph-duotone ph-columns" style="font-size:15px;color:var(--ic-search)"></i>Split
			</button>
			<button class="mode-btn" class:mode-btn--active={viewMode === 'read'}
				onclick={() => { viewMode = 'read'; }} type="button">
				<i class="ph-duotone ph-book-open-text" style="font-size:15px;color:var(--ic-doc)"></i>Read
			</button>
		</div>
		<div class="toolbar-divider"></div>
		<div class="doc-actions">
			<button class="doc-action-btn" type="button" title="Pin document">
				<i class="ph-duotone ph-push-pin" style="font-size:18px;color:var(--ic-pin)"></i>
			</button>
			<button class="doc-action-btn" type="button" title="Export">
				<i class="ph-duotone ph-export" style="font-size:18px;color:var(--ic-recent)"></i>
			</button>
			<button
				class="doc-action-btn doc-action-btn--save"
				class:doc-action-btn--saving={saving}
				disabled={saving}
				onclick={handleSave}
				type="button"
				title="Save (Ctrl+S)"
			>
				<i class="ph-duotone ph-floppy-disk" style="font-size:16px"></i>
				{saving ? 'Saving…' : 'Save'}
			</button>
		</div>
		{#if error}<p class="toolbar-error" role="alert">{error}</p>{/if}
	</div>

	<div class="workspace">
		{#if viewMode !== 'read'}
			<section class="pane pane--source" class:pane--half={viewMode === 'split'}>
				<Editor
					content={$activeDocument.content}
					title={$activeDocument.title}
					onContentChange={(v) => handleContentChange($activeDocument!.id, v)}
					onTitleChange={(v) => handleTitleChange($activeDocument!.id, v)}
					onCursorChange={handleCursorChange}
					onSave={handleSave}
					{saving}
				/>
			</section>
		{/if}
		{#if viewMode !== 'edit'}
			<section class="pane pane--preview" class:pane--half={viewMode === 'split'}>
				<div class="preview-header">
					<span class="preview-label">Preview</span>
				</div>
				<article class="preview-body thin-scroll">
					<div class="preview-content prose">
						{@html previewHtml}
					</div>
				</article>
			</section>
		{/if}
	</div>

{:else}
	<div class="home" aria-label="Welcome to ink-mark">
		<div class="home-gradient" aria-hidden="true"></div>
		<div class="home-mica" aria-hidden="true"></div>
		<div class="home-content">
			<div class="brand">
				<span class="brand-icon" aria-hidden="true">🖊️</span>
				<span class="brand-name">ink-mark</span>
			</div>
			<div class="hero">
				<h1 class="hero-title">Write without&nbsp;friction.</h1>
				<p class="hero-desc">
					A clean, distraction-free Markdown editor that saves your documents
					exactly where you want them — local files, your rules, no sync required.
				</p>
			</div>
			<div class="cta-row">
				<button class="cta-btn cta-btn--primary" onclick={handleNewDocument}
					disabled={loadingNew} type="button">
					{#if loadingNew}
						<i class="ph-duotone ph-circle-notch cta-icon cta-icon--spin" aria-hidden="true"></i>
					{:else}
						<i class="ph-duotone ph-note-pencil cta-icon" aria-hidden="true"></i>
					{/if}
					New Document
				</button>
				<button class="cta-btn cta-btn--secondary" onclick={handleOpenDocument}
					disabled={loadingOpen} type="button">
					{#if loadingOpen}
						<i class="ph-duotone ph-circle-notch cta-icon cta-icon--spin" aria-hidden="true"></i>
					{:else}
						<i class="ph-duotone ph-folder-open cta-icon" aria-hidden="true"></i>
					{/if}
					Open Document
				</button>
			</div>
			{#if error}
				<p class="home-error" role="alert">
					<i class="ph-duotone ph-warning-circle" aria-hidden="true"></i>
					{error}
				</p>
			{/if}
		</div>
	</div>
{/if}


<style>
	/* ── Doc toolbar ─────────────────────── */
	.doc-toolbar {
		display: flex; align-items: center; gap: 8px;
		height: 48px; padding: 0 12px; border-bottom: 1px solid var(--border);
		flex-shrink: 0; font-family: var(--font); font-size: 12.5px;
	}
	.breadcrumb { display: flex; align-items: center; gap: 6px; color: var(--ink-muted); min-width: 0; overflow: hidden; }
	.breadcrumb-cur { font-weight: 600; color: var(--ink); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 200px; }
	.save-status { display: flex; align-items: center; gap: 5px; font-size: 12px; color: var(--ink-subtle); white-space: nowrap; }
	.toolbar-spacer { flex: 1; }
	.toolbar-divider { width: 1px; height: 20px; background: var(--border); flex-shrink: 0; }
	.mode-switch { display: flex; align-items: center; padding: 2px; border-radius: 9px; background: var(--wash-hover); border: 1px solid var(--border); }
	.mode-btn {
		display: flex; align-items: center; gap: 5px; height: 28px; padding: 0 10px;
		border: none; border-radius: 7px; background: transparent; color: var(--ink-muted);
		font-family: var(--font); font-size: 12px; font-weight: 500;
		cursor: pointer; transition: background 100ms, color 100ms; white-space: nowrap;
	}
	.mode-btn:hover { color: var(--ink); }
	.mode-btn--active { background: var(--surface-panel); color: var(--ink); font-weight: 600; box-shadow: 0 1px 3px rgba(0,0,0,0.06); border: 1px solid var(--border); }
	.doc-actions { display: flex; align-items: center; gap: 2px; }
	.doc-action-btn { display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; border: none; border-radius: 7px; background: transparent; cursor: pointer; transition: background 100ms; }
	.doc-action-btn:hover { background: var(--wash-hover); }
	.doc-action-btn--save { width: auto; padding: 0 12px; gap: 5px; background: var(--accent); color: #fff; font-family: var(--font); font-size: 12.5px; font-weight: 600; transition: background 120ms, opacity 120ms; }
	.doc-action-btn--save:hover:not(:disabled) { background: var(--accent-hover); }
	.doc-action-btn--save:disabled, .doc-action-btn--saving { opacity: 0.60; cursor: not-allowed; }
	.toolbar-error { font-size: 11.5px; color: #C42B1C; margin: 0; white-space: nowrap; }

	/* ── Workspace ───────────────────────── */
	.workspace { flex: 1; display: flex; min-height: 0; overflow: hidden; }
	.pane { display: flex; flex-direction: column; min-width: 0; flex: 1; overflow: hidden; }
	.pane--half { flex: 1 1 50%; }
	.pane--preview { border-left: 1px solid var(--border); }
	.preview-header { display: flex; align-items: center; height: 38px; padding: 0 16px; border-bottom: 1px solid var(--border); flex-shrink: 0; }
	.preview-label { font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: var(--ink-subtle); font-family: var(--font); }
	.preview-body { flex: 1; overflow-y: auto; }
	.preview-content { max-width: 680px; margin: 0 auto; padding: 32px 40px; }
	/* Prose (rendered markdown) */
	.prose { font-family: var(--font); font-size: 15px; line-height: 1.70; color: var(--ink); }
	.prose :global(h1) { font-size: 26px; font-weight: 700; letter-spacing: -0.025em; margin: 0 0 16px; }
	.prose :global(h2) { font-size: 20px; font-weight: 700; margin: 28px 0 12px; }
	.prose :global(h3) { font-size: 16px; font-weight: 600; margin: 24px 0 8px; }
	.prose :global(p)  { margin: 0 0 14px; }
	.prose :global(ul), .prose :global(ol) { padding-left: 24px; margin: 0 0 14px; }
	.prose :global(li) { margin-bottom: 4px; }
	.prose :global(blockquote) { border-left: 3px solid var(--accent); padding-left: 16px; margin: 16px 0; color: var(--ink-muted); font-style: italic; }
	.prose :global(code) { background: rgba(15,108,189,0.07); color: var(--accent); border-radius: 4px; padding: 1px 6px; font-size: 0.88em; }
	.prose :global(pre) { background: #f0f4f9; border-radius: 8px; padding: 16px; overflow-x: auto; margin: 0 0 16px; }
	.prose :global(pre code) { background: none; padding: 0; font-size: 13px; }
	.prose :global(a) { color: var(--accent); text-underline-offset: 2px; }
	.prose :global(hr) { border: none; border-top: 1px solid var(--border); margin: 24px 0; }
	.prose :global(strong) { font-weight: 700; }
	.prose :global(table) { width: 100%; border-collapse: collapse; margin: 0 0 16px; font-size: 14px; }
	.prose :global(th) { text-align: left; font-weight: 600; padding: 8px 12px; border-bottom: 2px solid var(--border-strong); color: var(--ink-muted); }
	.prose :global(td) { padding: 7px 12px; border-bottom: 1px solid var(--border); }

	/* ── HOME PAGE ────────────────────────── */
	.home { position: relative; display: flex; align-items: center; justify-content: center; flex: 1; overflow: hidden; user-select: none; }
	.home-gradient {
		position: absolute; inset: 0;
		background: radial-gradient(ellipse 80% 60% at 20% 10%, rgba(196,220,248,0.55) 0%, transparent 65%),
			radial-gradient(ellipse 60% 50% at 85% 80%, rgba(180,214,247,0.45) 0%, transparent 60%),
			linear-gradient(155deg, #EEF5FD 0%, #F5F9FF 40%, #EDF4FD 70%, #F0F6FE 100%);
		z-index: 0;
	}
	.home-mica { position: absolute; inset: 0; backdrop-filter: blur(36px) saturate(160%) brightness(1.04); -webkit-backdrop-filter: blur(36px) saturate(160%) brightness(1.04); background: rgba(255,255,255,0.58); z-index: 1; }
	.home-content {
		position: relative; z-index: 2; display: flex; flex-direction: column; align-items: center; text-align: center;
		background: rgba(255,255,255,0.72); border: 1px solid rgba(255,255,255,0.90); border-radius: 20px;
		padding: 52px 60px 48px; max-width: 520px; width: calc(100% - 48px);
		box-shadow: 0 2px 1px rgba(0,0,0,0.03), 0 6px 16px rgba(15,108,189,0.06), 0 20px 48px rgba(15,108,189,0.07), inset 0 1px 0 rgba(255,255,255,0.95);
	}
	.brand { display: flex; align-items: center; gap: 8px; margin-bottom: 36px; }
	.brand-icon { font-size: 22px; line-height: 1; }
	.brand-name { font-family: 'Plus Jakarta Sans', var(--font); font-size: 15px; font-weight: 700; color: #1A1A1A; letter-spacing: -0.02em; }
	.hero { display: flex; flex-direction: column; align-items: center; gap: 16px; margin-bottom: 40px; }
	.hero-title { font-family: 'Plus Jakarta Sans', var(--font); font-size: 36px; font-weight: 800; color: #0D1117; letter-spacing: -0.035em; line-height: 1.12; margin: 0; }
	.hero-desc { font-family: var(--font); font-size: 14.5px; line-height: 1.70; color: #505668; margin: 0; max-width: 360px; }
	.cta-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; justify-content: center; }
	.cta-btn { display: inline-flex; align-items: center; gap: 8px; padding: 10px 22px; border-radius: 8px; font-family: var(--font); font-size: 13.5px; font-weight: 600; cursor: pointer; transition: background 130ms, box-shadow 130ms, transform 80ms; outline: none; white-space: nowrap; }
	.cta-btn--primary { background: var(--accent); color: #fff; border: 1px solid var(--accent); box-shadow: 0 1px 3px rgba(15,108,189,0.20), 0 4px 12px rgba(15,108,189,0.18); }
	.cta-btn--primary:hover:not(:disabled) { background: var(--accent-hover); border-color: var(--accent-hover); }
	.cta-btn--primary:active:not(:disabled) { background: #0B5396; transform: scale(0.97); }
	.cta-btn--secondary { background: rgba(15,108,189,0.06); color: var(--accent); border: 1.5px solid rgba(15,108,189,0.30); }
	.cta-btn--secondary:hover:not(:disabled) { background: rgba(15,108,189,0.11); border-color: rgba(15,108,189,0.50); }
	.cta-btn--secondary:active:not(:disabled) { background: rgba(15,108,189,0.16); transform: scale(0.97); }
	.cta-btn:disabled { opacity: 0.55; cursor: not-allowed; }
	.cta-icon { font-size: 17px; flex-shrink: 0; }
	@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
	.cta-icon--spin { animation: spin 900ms linear infinite; }
	.home-error { display: flex; align-items: center; gap: 6px; margin: 20px 0 0; font-family: var(--font); font-size: 12px; color: #C42B1C; }
</style>


