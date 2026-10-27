<script lang="ts">
	import {
		newDocumentTrigger,
		activeDocumentId,
		allDocuments,
		allProjects,
		openDocuments,
		openInTab,
		dirtyDocumentIds,
		markClean
	} from '#lib/stores';
	import {
		ListDocuments, ListProjects, ListFavorited,
		LoadDocument, NewProject, OpenFolderAsProject,
		FavoriteDocument, UnfavoriteDocument
	} from '#lib/wailsjs/go/main/App';
	import type { Document, Project } from '#lib/types';
	import { get } from 'svelte/store';

	let createMenuOpen = $state(false);

	function toggleCreateMenu(e: MouseEvent) {
		e.stopPropagation();
		createMenuOpen = !createMenuOpen;
	}
	function closeCreateMenu() { createMenuOpen = false; }

	// Load documents, projects and favorites on mount.
	$effect(() => {
		Promise.all([ListDocuments(), ListProjects(), ListFavorited()])
			.then(([docs, projects, favs]) => {
				allDocuments.set(docs);
				allProjects.set(projects);
				// Merge favorited flag into allDocuments from the favorites list.
				const favIds = new Set(favs.map((d) => d.id));
				allDocuments.update((all) =>
					all.map((d) => ({ ...d, favorited: favIds.has(d.id) }))
				);
			})
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

	async function handleNewProject(): Promise<void> {
		closeCreateMenu();
		try {
			const p = await NewProject();
			if (p.id) allProjects.update((ps) => {
				const exists = ps.some((x) => x.id === p.id);
				return exists ? ps : [p as Project, ...ps];
			});
		} catch (e) {
			console.error('Failed to create project:', e);
		}
	}

	async function handleOpenFolderAsProject(): Promise<void> {
		closeCreateMenu();
		try {
			const p = await OpenFolderAsProject();
			if (p.id) {
				allProjects.update((ps) => {
					const exists = ps.some((x) => x.id === p.id);
					return exists ? ps : [p as Project, ...ps];
				});
				// Register all project docs into allDocuments too.
				if (p.documents?.length) {
					allDocuments.update((docs) => {
						const ids = new Set(docs.map((d) => d.id));
						const newDocs = p.documents.filter((d) => !ids.has(d.id));
						return [...newDocs, ...docs];
					});
				}
			}
		} catch (e) {
			console.error('Failed to open folder as project:', e);
		}
	}

	async function handleToggleFavorite(doc: Document, e: MouseEvent): Promise<void> {
		e.stopPropagation();
		try {
			if (doc.favorited) {
				await UnfavoriteDocument(doc.id);
			} else {
				await FavoriteDocument(doc.id);
			}
			// Update in allDocuments.
			allDocuments.update((docs) =>
				docs.map((d) => d.id === doc.id ? { ...d, favorited: !doc.favorited } : d)
			);
		} catch (e) {
			console.error('Failed to toggle favorite:', e);
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
					<button class="create-item" type="button" role="menuitem" onclick={handleNewProject}>
						<i class="ph-duotone ph-folder-simple-plus create-item-icon" style="color:var(--ic-project)"></i>
						<span class="create-item-body">
							<span class="create-item-title">New project</span>
							<span class="create-item-sub">A folder of related docs</span>
						</span>
					</button>
					<div class="create-divider"></div>
					<button class="create-item" type="button" role="menuitem">
						<i class="ph-duotone ph-file-arrow-up create-item-icon" style="color:var(--ic-recent)"></i>
						<span class="create-item-body">Open file</span>
					</button>
					<button class="create-item" type="button" role="menuitem" onclick={handleOpenFolderAsProject}>
						<i class="ph-duotone ph-folder-open create-item-icon" style="color:var(--ic-project)"></i>
						<span class="create-item-body">Open folder as project</span>
					</button>
				</div>
			{/if}
		</div>
	</div>


	<nav class="sb-nav thin-scroll">

		<!-- Favorites -->
		<section class="sb-section">
			<div class="sb-section-header">
				<span class="sb-section-title">Favorites</span>
			</div>
			<div class="sb-tree">
				{#each $allDocuments.filter(d => d.favorited) as doc (doc.id)}
					<button
						class="sb-file-link"
						class:sb-file-link--active={$activeDocumentId === doc.id}
						type="button"
						onclick={() => handleDocClick(doc.id)}
					>
						<i class="ph-duotone ph-star sb-file-icon" style="color:var(--ic-pin)"></i>
						<span class="sb-file-name">{doc.title}</span>
						{#if $dirtyDocumentIds.has(doc.id)}
							<span class="sb-unsaved-dot" title="Unsaved changes"></span>
						{/if}
					</button>
				{:else}
					<p class="sb-empty">No favorites yet.</p>
				{/each}
			</div>
		</section>

		<!-- Documents -->
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
						{#if $dirtyDocumentIds.has(doc.id)}
							<span class="sb-unsaved-dot" title="Unsaved changes"></span>
						{/if}
						<span
						class="sb-fav-btn"
						role="button"
						tabindex="0"
						title={doc.favorited ? 'Remove from favorites' : 'Add to favorites'}
						onclick={(e) => handleToggleFavorite(doc, e)}
							onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleToggleFavorite(doc, e as unknown as MouseEvent); }}
					>
							<i class="ph-duotone ph-star" style="font-size:13px;color:{doc.favorited ? 'var(--ic-pin)' : 'var(--ink-subtle)'}"></i>
						</span>
					</button>
				{:else}
					<p class="sb-empty">No documents yet.</p>
				{/each}

			</div>
		</section>

		<!-- Projects -->
		<section class="sb-section">
			<div class="sb-section-header">
				<span class="sb-section-title">Projects</span>
				<button class="sb-section-add" type="button" title="New project"
					onclick={handleNewProject}>
					<i class="ph-duotone ph-plus" style="font-size:14px;color:var(--ic-create)"></i>
				</button>
			</div>
			{#each $allProjects as project (project.id)}
				<details class="sb-project">
					<summary class="sb-project-summary">
						<i class="ph-duotone ph-caret-down sb-project-caret"></i>
						<i class="ph-duotone ph-books sb-project-icon" style="color:var(--ic-project)"></i>
						<span class="sb-project-name">{project.name}</span>
						<span class="sb-project-count">{project.documents.length}</span>
					</summary>
					<div class="sb-project-children">
						{#each project.documents as doc (doc.id)}
							<button
								class="sb-file-link"
								class:sb-file-link--active={$activeDocumentId === doc.id}
								type="button"
								onclick={() => handleDocClick(doc.id)}
							>
								<i class="ph-duotone ph-file-text sb-file-icon" style="color:var(--ic-doc)"></i>
								<span class="sb-file-name">{doc.title}</span>
								{#if $dirtyDocumentIds.has(doc.id)}
									<span class="sb-unsaved-dot" title="Unsaved changes"></span>
								{/if}
							</button>
						{:else}
							<p class="sb-empty">No documents.</p>
						{/each}
					</div>
				</details>
			{:else}
				<p class="sb-empty">No projects yet.</p>
			{/each}
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
	.sb-nav-icon { font-size: 20px; flex-shrink: 0; }
	.sb-nav-label { flex: 1; font-size: 13px; }
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
	.sb-unsaved-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--ic-doc); flex-shrink: 0; }
	/* .sb-indent used when subfolder nesting is wired (backlog item 6) */
	:global(.sb-indent) { padding-left: 10px; }
	.sb-empty { font-size: 12px; color: var(--ink-subtle); font-family: var(--font); padding: 4px 8px; margin: 0; }
	.sb-fav-btn { display: none; align-items: center; justify-content: center; width: 18px; height: 18px; border: none; border-radius: 3px; background: transparent; cursor: pointer; flex-shrink: 0; padding: 0; }
	.sb-file-link:hover .sb-fav-btn { display: flex; }
	.sb-file-link--active .sb-fav-btn { display: flex; }

	/* Projects & folders — folder selectors used in backlog item 6 */
	.sb-project { margin-bottom: 2px; }
	.sb-project summary { list-style: none; }
	.sb-project summary::-webkit-details-marker { display: none; }
	.sb-project-summary { display: flex; align-items: center; gap: 5px; height: 34px; padding: 0 4px; border-radius: 8px; cursor: pointer; transition: background 100ms; user-select: none; }
	.sb-project-summary:hover { background: var(--wash-hover); }
	.sb-project-caret { font-size: 14px; color: var(--ink-subtle); flex-shrink: 0; transition: transform 150ms; }
	.sb-project:not([open]) .sb-project-caret { transform: rotate(-90deg); }
	.sb-project-icon { font-size: 22px; flex-shrink: 0; }
	.sb-project-name { flex: 1; font-size: 13px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.sb-project-count { font-size: 11px; color: var(--ink-subtle); padding-right: 4px; }
	.sb-project-children { margin-left: 16px; padding-left: 8px; border-left: 1px solid var(--border); display: flex; flex-direction: column; gap: 1px; }
	:global(.sb-folder) { margin-bottom: 1px; }
	:global(.sb-folder-summary) { display: flex; align-items: center; gap: 5px; height: 29px; padding: 0 4px; border-radius: 7px; cursor: pointer; transition: background 100ms; user-select: none; }
	:global(.sb-folder-summary:hover) { background: var(--wash-hover); }
	:global(.sb-folder-caret) { font-size: 13px; color: var(--ink-subtle); flex-shrink: 0; transition: transform 150ms; }
	:global(.sb-folder:not([open]) .sb-folder-caret) { transform: rotate(-90deg); }
	:global(.sb-folder-icon) { font-size: 18px; flex-shrink: 0; }
	:global(.sb-folder-name) { flex: 1; font-size: 12.5px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	:global(.sb-folder-children) { margin-left: 13px; padding-left: 6px; border-left: 1px solid var(--border); display: flex; flex-direction: column; gap: 1px; }

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


