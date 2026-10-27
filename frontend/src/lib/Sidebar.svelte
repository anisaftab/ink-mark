<script lang="ts">
	import {
		newDocumentTrigger,
		activeDocumentId,
		allDocuments,
		openDocuments,
		openInTab
	} from '#lib/stores';
	import { ListDocuments, LoadDocument } from '#lib/wailsjs/go/main/App';
	import type { Document } from '#lib/types';
	import { get } from 'svelte/store';

	let createMenuOpen = $state(false);

	function toggleCreateMenu(e: MouseEvent) {
		e.stopPropagation();
		createMenuOpen = !createMenuOpen;
	}
	function closeCreateMenu() { createMenuOpen = false; }

	// Load the document list once on mount.
	$effect(() => {
		ListDocuments()
			.then((docs) => allDocuments.set(docs))
			.catch(console.error);
	});

	async function handleDocClick(id: string): Promise<void> {
		const already = get(openDocuments).find((d) => d.id === id);
		if (already) { activeDocumentId.set(id); return; }
		try {
			const dto = await LoadDocument(id);
			if (dto.id) openInTab(dto as Document);
		} catch (e) {
			console.error('Failed to load document:', e);
		}
	}
</script>


<svelte:window onclick={closeCreateMenu} />

<aside class="sidebar glass-sidebar">
	<div class="sb-brand">
		<i class="ph-duotone ph-notebook sb-brand-icon" style="color:var(--ic-doc)"></i>
		<span class="sb-brand-name">ink-mark</span>
	</div>

	<div class="sb-actions">
		<button class="sb-search" type="button">
			<i class="ph-duotone ph-magnifying-glass sb-search-icon" style="color:var(--ic-search)"></i>
			<span class="sb-search-label">Search everything</span>
			<kbd class="sb-kbd">Ctrl K</kbd>
		</button>
		<div class="sb-create-row">
			<button class="sb-create-main" type="button"
				onclick={() => newDocumentTrigger.update(n => n + 1)}>
				<i class="ph-duotone ph-file-plus sb-create-icon"></i>
				New document
			</button>
			<button class="sb-create-caret" type="button" title="More create options"
				onclick={toggleCreateMenu}>
				<i class="ph-duotone ph-caret-down" style="font-size:13px"></i>
			</button>
			{#if createMenuOpen}
				<div class="create-dropdown glass-content" role="menu">
					<button class="create-item" type="button" role="menuitem"
						onclick={() => { newDocumentTrigger.update(n => n + 1); closeCreateMenu(); }}>
						<i class="ph-duotone ph-file-text create-item-icon" style="color:var(--ic-doc)"></i>
						<span class="create-item-body">
							<span class="create-item-title">New document</span>
							<span class="create-item-sub">A standalone .md file</span>
						</span>
						<kbd class="create-item-kbd">Ctrl N</kbd>
					</button>
					<button class="create-item" type="button" role="menuitem">
						<i class="ph-duotone ph-folder-simple-plus create-item-icon" style="color:var(--ic-project)"></i>
						<span class="create-item-body">
							<span class="create-item-title">New project</span>
							<span class="create-item-sub">A folder of related docs</span>
						</span>
						<kbd class="create-item-kbd">Ctrl ⇧N</kbd>
					</button>
					<div class="create-divider"></div>
					<button class="create-item" type="button" role="menuitem">
						<i class="ph-duotone ph-file-arrow-up create-item-icon" style="color:var(--ic-recent)"></i>
						<span class="create-item-body">Open file</span>
						<kbd class="create-item-kbd">Ctrl O</kbd>
					</button>
					<button class="create-item" type="button" role="menuitem">
						<i class="ph-duotone ph-folder-open create-item-icon" style="color:var(--ic-project)"></i>
						<span class="create-item-body">Open folder as project</span>
					</button>
				</div>
			{/if}
		</div>
	</div>

	<nav class="sb-nav thin-scroll">
		<div class="sb-quick">
			<a href="#" class="sb-nav-link">
				<i class="ph-duotone ph-clock-counter-clockwise sb-nav-icon" style="color:var(--ic-recent)"></i>
				<span class="sb-nav-label">Recent</span>
			</a>
			<a href="#" class="sb-nav-link">
				<i class="ph-duotone ph-push-pin sb-nav-icon" style="color:var(--ic-pin)"></i>
				<span class="sb-nav-label">Pinned</span>
				<span class="sb-nav-badge">3</span>
			</a>
		</div>
		<section class="sb-section">
			<div class="sb-section-header">
				<span class="sb-section-title">Documents</span>
				<button class="sb-section-add" type="button" title="New document"
					onclick={() => newDocumentTrigger.update(n => n + 1)}>
					<i class="ph-duotone ph-plus" style="font-size:14px;color:var(--ic-create)"></i>
				</button>
			</div>
			<div class="sb-tree">
				{#each $allDocuments as doc (doc.id)}
					<button
						class="sb-file-link"
						class:sb-file-link--active={$activeDocumentId === doc.id}
						type="button"
						onclick={() => handleDocClick(doc.id)}
					>
						<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
						<span class="sb-file-name">{doc.title}</span>
					</button>
				{:else}
					<p class="sb-empty">No documents yet.</p>
				{/each}
			</div>
		</section>

		<!-- Projects (placeholder — wired in backlog item 6) -->
		<section class="sb-section">
			<div class="sb-section-header">
				<span class="sb-section-title">Projects</span>
				<button class="sb-section-add" type="button" title="New project">
					<i class="ph-duotone ph-plus" style="font-size:14px;color:var(--ic-create)"></i>
				</button>
			</div>
			<details class="sb-project" open>
				<summary class="sb-project-summary">
					<i class="ph-duotone ph-caret-down sb-project-caret"></i>
					<i class="ph-duotone ph-books sb-project-icon" style="color:var(--ic-project)"></i>
					<span class="sb-project-name">Atlas Handbook</span>
					<span class="sb-project-count">9</span>
				</summary>
				<div class="sb-project-children">
					<details class="sb-folder" open>
						<summary class="sb-folder-summary">
							<i class="ph-duotone ph-caret-down sb-folder-caret"></i>
							<i class="ph-duotone ph-folder-open sb-folder-icon" style="color:var(--ic-project)"></i>
							<span class="sb-folder-name">research</span>
						</summary>
						<div class="sb-folder-children">
							<a href="#" class="sb-file-link sb-indent">
								<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
								<span class="sb-file-name">interview-notes</span>
							</a>
							<a href="#" class="sb-file-link sb-indent">
								<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
								<span class="sb-file-name">q3-synthesis</span>
								<span class="sb-unsaved-dot" title="Unsaved changes"></span>
							</a>
							<a href="#" class="sb-file-link sb-indent">
								<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
								<span class="sb-file-name">personas</span>
							</a>
						</div>
					</details>
					<a href="#" class="sb-file-link">
						<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
						<span class="sb-file-name">README</span>
					</a>
				</div>
			</details>
			<details class="sb-project">
				<summary class="sb-project-summary">
					<i class="ph-duotone ph-caret-down sb-project-caret"></i>
					<i class="ph-duotone ph-books sb-project-icon" style="color:var(--ic-project)"></i>
					<span class="sb-project-name">Kestrel Blog</span>
					<span class="sb-project-count">12</span>
				</summary>
				<div class="sb-project-children">
					<a href="#" class="sb-file-link">
						<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
						<span class="sb-file-name">drafts-index</span>
					</a>
				</div>
			</details>
			<details class="sb-project">
				<summary class="sb-project-summary">
					<i class="ph-duotone ph-caret-down sb-project-caret"></i>
					<i class="ph-duotone ph-books sb-project-icon" style="color:var(--ic-project)"></i>
					<span class="sb-project-name">Thesis Draft</span>
					<span class="sb-project-count">27</span>
				</summary>
				<div class="sb-project-children">
					<a href="#" class="sb-file-link">
						<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
						<span class="sb-file-name">chapter-01</span>
					</a>
				</div>
			</details>
		</section>
	</nav>

	<div class="sb-footer">
		<div class="sb-vault">
			<i class="ph-duotone ph-lock-simple sb-vault-icon" style="color:var(--ic-secure)"></i>
			<div class="sb-vault-info">
				<span class="sb-vault-label">Private - On this device</span>
				<span class="sb-vault-path">~/Documents/ink-mark</span>
			</div>
		</div>
		<div class="sb-footer-row">
			<button class="sb-footer-btn sb-footer-btn--wide" type="button">
				<i class="ph-duotone ph-gear-six sb-nav-icon" style="color:var(--ic-settings)"></i>
				<span class="sb-nav-label">Settings</span>
				<kbd class="sb-kbd">Ctrl ,</kbd>
			</button>
			<button class="sb-footer-icon-btn" type="button" title="Toggle theme">
				<i class="ph-duotone ph-moon-stars" style="font-size:18px;color:var(--ic-outline)"></i>
			</button>
			<button class="sb-footer-icon-btn" type="button" title="Help & shortcuts">
				<i class="ph-duotone ph-question" style="font-size:18px;color:var(--ic-recent)"></i>
			</button>
		</div>
	</div>
</aside>


<style>
	.sidebar { display: flex; flex-direction: column; height: 100%; border-right: 1px solid var(--border); overflow: hidden; }
	.sb-brand { display: flex; align-items: center; gap: 8px; height: 48px; padding: 0 16px; flex-shrink: 0; }
	.sb-brand-icon { font-size: 22px; line-height: 1; }
	.sb-brand-name { font-size: 15px; font-weight: 700; color: var(--ink); letter-spacing: -0.02em; }
	.sb-actions { padding: 0 12px 12px; display: flex; flex-direction: column; gap: 8px; flex-shrink: 0; }
	.sb-search { display: flex; align-items: center; gap: 8px; width: 100%; height: 36px; padding: 0 10px; border: 1px solid var(--border); border-radius: 9px; background: var(--wash-hover); color: var(--ink-subtle); font-family: var(--font); font-size: 12.5px; cursor: pointer; transition: border-color 120ms, background 120ms; text-align: left; }
	.sb-search:hover { border-color: var(--border-strong); background: var(--wash-selected); }
	.sb-search-icon { font-size: 18px; flex-shrink: 0; }
	.sb-search-label { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.sb-kbd { font-family: var(--font); font-size: 11px; font-weight: 500; color: var(--ink-muted); padding: 1px 5px; border: 1px solid var(--border); border-radius: 4px; white-space: nowrap; flex-shrink: 0; }
	.sb-create-row { position: relative; display: flex; }
	.sb-create-main { display: flex; align-items: center; justify-content: center; gap: 7px; flex: 1; height: 36px; padding: 0 14px; border: none; border-radius: 9px 0 0 9px; background: var(--ink); color: var(--surface-panel); font-family: var(--font); font-size: 12.5px; font-weight: 600; cursor: pointer; transition: opacity 120ms; }
	.sb-create-main:hover { opacity: 0.88; }
	.sb-create-icon { font-size: 18px; }
	.sb-create-caret { display: flex; align-items: center; justify-content: center; width: 32px; height: 36px; border: none; border-left: 1px solid rgba(253,252,250,0.15); border-radius: 0 9px 9px 0; background: var(--ink); color: var(--surface-panel); cursor: pointer; transition: opacity 120ms; }
	.sb-create-caret:hover { opacity: 0.88; }
	.create-dropdown { position: absolute; z-index: 50; top: calc(100% + 6px); left: 0; right: 0; border: 1px solid var(--border-strong); border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.04), 0 12px 32px rgba(0,0,0,0.10); padding: 6px; }
	.create-item { display: flex; align-items: center; gap: 10px; width: 100%; padding: 7px 10px; border: none; border-radius: 8px; background: transparent; color: var(--ink); font-family: var(--font); font-size: 12.5px; cursor: pointer; text-align: left; transition: background 100ms; }
	.create-item:hover { background: var(--wash-hover); }
	.create-item-icon { font-size: 19px; flex-shrink: 0; }
	.create-item-body { display: flex; flex-direction: column; flex: 1; min-width: 0; }
	.create-item-title { font-weight: 600; }
	.create-item-sub { font-size: 11px; color: var(--ink-muted); }
	.create-item-kbd { font-family: var(--font); font-size: 11px; color: var(--ink-subtle); flex-shrink: 0; }
	.create-divider { height: 1px; background: var(--border); margin: 5px 0; }
	.sb-nav { flex: 1; overflow-y: auto; padding: 4px 8px 12px; display: flex; flex-direction: column; }
	.sb-quick { display: flex; flex-direction: column; gap: 2px; margin-bottom: 16px; }
	.sb-nav-link { display: flex; align-items: center; gap: 9px; height: 34px; padding: 0 8px; border-radius: 8px; color: var(--ink); text-decoration: none; font-weight: 500; transition: background 100ms; }
	.sb-nav-link:hover { background: var(--wash-hover); }
	.sb-nav-icon { font-size: 20px; flex-shrink: 0; }
	.sb-nav-label { flex: 1; font-size: 13px; }
	.sb-nav-badge { font-size: 11px; color: var(--ink-subtle); }
	.sb-section { margin-bottom: 16px; }
	.sb-section-header { display: flex; align-items: center; justify-content: space-between; height: 26px; padding: 0 8px; margin-bottom: 2px; }
	.sb-section-title { font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: var(--ink-subtle); }
	.sb-section-add { display: flex; align-items: center; justify-content: center; width: 22px; height: 22px; border: none; border-radius: 6px; background: transparent; cursor: pointer; transition: background 100ms; }
	.sb-section-add:hover { background: var(--wash-selected); }
	.sb-tree { display: flex; flex-direction: column; gap: 1px; }
	.sb-file-link { display: flex; align-items: center; gap: 8px; width: 100%; height: 30px; padding: 0 8px; border: none; border-radius: 7px; background: transparent; color: var(--ink); text-decoration: none; text-align: left; font-family: var(--font); cursor: pointer; transition: background 100ms; }
	.sb-file-link:hover { background: var(--wash-hover); }
	.sb-file-link--active { background: var(--wash-selected); font-weight: 600; }
	.sb-file-icon { font-size: 18px; flex-shrink: 0; }
	.sb-file-name { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 12.5px; }
	.sb-pin-icon { font-size: 14px; flex-shrink: 0; }
	.sb-unsaved-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--ic-doc); flex-shrink: 0; }
	.sb-indent { padding-left: 10px; }
	.sb-empty { font-size: 12px; color: var(--ink-subtle); font-family: var(--font); padding: 4px 8px; margin: 0; }

	/* Projects & folders */
	.sb-project { margin-bottom: 2px; }
	.sb-project summary, .sb-folder summary { list-style: none; }
	.sb-project summary::-webkit-details-marker, .sb-folder summary::-webkit-details-marker { display: none; }
	.sb-project-summary { display: flex; align-items: center; gap: 5px; height: 34px; padding: 0 4px; border-radius: 8px; cursor: pointer; transition: background 100ms; user-select: none; }
	.sb-project-summary:hover { background: var(--wash-hover); }
	.sb-project-caret { font-size: 14px; color: var(--ink-subtle); flex-shrink: 0; transition: transform 150ms; }
	.sb-project:not([open]) .sb-project-caret { transform: rotate(-90deg); }
	.sb-project-icon { font-size: 22px; flex-shrink: 0; }
	.sb-project-name { flex: 1; font-size: 13px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.sb-project-count { font-size: 11px; color: var(--ink-subtle); padding-right: 4px; }
	.sb-project-children { margin-left: 16px; padding-left: 8px; border-left: 1px solid var(--border); display: flex; flex-direction: column; gap: 1px; }
	.sb-folder { margin-bottom: 1px; }
	.sb-folder-summary { display: flex; align-items: center; gap: 5px; height: 29px; padding: 0 4px; border-radius: 7px; cursor: pointer; transition: background 100ms; user-select: none; }
	.sb-folder-summary:hover { background: var(--wash-hover); }
	.sb-folder-caret { font-size: 13px; color: var(--ink-subtle); flex-shrink: 0; transition: transform 150ms; }
	.sb-folder:not([open]) .sb-folder-caret { transform: rotate(-90deg); }
	.sb-folder-icon { font-size: 18px; flex-shrink: 0; }
	.sb-folder-name { flex: 1; font-size: 12.5px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.sb-folder-children { margin-left: 13px; padding-left: 6px; border-left: 1px solid var(--border); display: flex; flex-direction: column; gap: 1px; }

	/* Footer */
	.sb-footer { border-top: 1px solid var(--border); padding: 8px; flex-shrink: 0; display: flex; flex-direction: column; gap: 4px; }
	.sb-vault { display: flex; align-items: center; gap: 9px; padding: 6px 8px; border-radius: 8px; }
	.sb-vault-icon { font-size: 19px; flex-shrink: 0; }
	.sb-vault-info { display: flex; flex-direction: column; min-width: 0; }
	.sb-vault-label { font-size: 12px; font-weight: 600; color: var(--ink); }
	.sb-vault-path { font-size: 11px; color: var(--ink-subtle); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.sb-footer-row { display: flex; align-items: center; gap: 2px; }
	.sb-footer-btn { display: flex; align-items: center; gap: 8px; height: 34px; padding: 0 8px; border: none; border-radius: 8px; background: transparent; color: var(--ink); font-family: var(--font); font-size: 13px; font-weight: 500; cursor: pointer; transition: background 100ms; }
	.sb-footer-btn:hover { background: var(--wash-hover); }
	.sb-footer-btn--wide { flex: 1; }
	.sb-footer-icon-btn { display: flex; align-items: center; justify-content: center; width: 34px; height: 34px; border: none; border-radius: 8px; background: transparent; cursor: pointer; transition: background 100ms; }
	.sb-footer-icon-btn:hover { background: var(--wash-hover); }
</style>


