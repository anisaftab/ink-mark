<script lang="ts">
	import './layout.css';
	import '#lib/devlog';
	import favicon from '#lib/assets/favicon.svg';
	import Sidebar from '#lib/Sidebar.svelte';
	import { openDocuments, activeDocumentId, newDocumentTrigger, closeTab } from '#lib/stores';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="backdrop">
	<div class="shell">
		<Sidebar />
		<main class="content glass-content">
			<div class="titlebar">
				<div class="tabs">
					{#each $openDocuments as doc (doc.id)}
						<button
							class="tab"
							class:tab--active={$activeDocumentId === doc.id}
							type="button"
							onclick={() => activeDocumentId.set(doc.id)}
						>
							<i class="ph-duotone ph-file-text tab-icon" style="color:var(--ic-doc)"></i>
							<span class="tab-name">{doc.title || 'Untitled'}</span>
							<span
								class="tab-close"
								class:tab-close--hidden={$activeDocumentId !== doc.id}
								role="button"
								tabindex="0"
								aria-label="Close tab"
								onclick={(e) => { e.stopPropagation(); closeTab(doc.id); }}
								onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.stopPropagation(); closeTab(doc.id); } }}
							>
								<i class="ph-duotone ph-x" style="font-size:13px;color:var(--ink-subtle)"></i>
							</span>
						</button>
					{/each}
					<button class="tab-new" type="button" title="New tab"
						onclick={() => newDocumentTrigger.update(n => n + 1)}>
						<i class="ph-duotone ph-plus" style="font-size:18px;color:var(--ic-create)"></i>
					</button>
				</div>
				<div class="titlebar-spacer"></div>
				<div class="window-controls">
					<button class="wc-btn" type="button" title="Minimize">
						<i class="ph-duotone ph-minus" style="font-size:14px;color:var(--ink-muted)"></i>
					</button>
					<button class="wc-btn" type="button" title="Maximize">
						<i class="ph-duotone ph-square" style="font-size:14px;color:var(--ink-muted)"></i>
					</button>
					<button class="wc-btn wc-btn--close" type="button" title="Close">
						<i class="ph-duotone ph-x" style="font-size:14px;color:var(--ink-muted)"></i>
					</button>
				</div>
			</div>

			{@render children()}

			<footer class="statusbar">
				<span class="status-item">Ln 1, Col 1</span>
				<span class="status-item">0 words</span>
				<div class="status-spacer"></div>
				<span class="status-item">UTF-8</span>
				<span class="status-item">Markdown</span>
				<span class="status-item status-item--secure">
					<i class="ph-duotone ph-shield-check" style="font-size:14px;color:var(--ic-secure)"></i>
					Stored locally
				</span>
			</footer>
		</main>
	</div>
</div>

<style>
	.backdrop { position: fixed; inset: 0; background: var(--surface-base); font-family: var(--font); color: var(--ink); font-size: 13px; }
	.shell { display: grid; grid-template-columns: 272px 1fr; height: 100vh; overflow: hidden; }
	.content { display: flex; flex-direction: column; height: 100%; overflow: hidden; min-width: 0; }
	.titlebar { display: flex; align-items: flex-end; height: 48px; flex-shrink: 0; border-bottom: 1px solid var(--border); padding-left: 8px; }
	.tabs { display: flex; align-items: flex-end; gap: 2px; overflow-x: auto; min-width: 0; }
	.tabs::-webkit-scrollbar { display: none; }
	.tab { display: flex; align-items: center; gap: 7px; height: 36px; padding: 0 10px; border: none; border-radius: 8px 8px 0 0; background: transparent; color: var(--ink-muted); font-family: var(--font); font-size: 12.5px; cursor: pointer; transition: background 100ms, color 100ms; white-space: nowrap; }
	.tab:hover { background: var(--wash-hover); color: var(--ink); }
	.tab--active { background: var(--surface-base); color: var(--ink); font-weight: 600; border: 1px solid var(--border); border-bottom: none; margin-bottom: -1px; }
	.tab-icon { font-size: 15px; flex-shrink: 0; }
	.tab-name { max-width: 140px; overflow: hidden; text-overflow: ellipsis; }
	.tab-close { display: flex; align-items: center; justify-content: center; width: 18px; height: 18px; border-radius: 4px; flex-shrink: 0; transition: background 100ms; }
	.tab-close:hover { background: var(--wash-selected); }
	.tab-close--hidden { opacity: 0; }
	.tab:hover .tab-close--hidden { opacity: 1; }
	.tab-new { display: flex; align-items: center; justify-content: center; width: 34px; height: 34px; border: none; border-radius: 8px; background: transparent; cursor: pointer; margin-bottom: 2px; transition: background 100ms; }
	.tab-new:hover { background: var(--wash-hover); }
	.titlebar-spacer { flex: 1; }
	.window-controls { display: flex; align-self: flex-start; height: 48px; }
	.wc-btn { display: flex; align-items: center; justify-content: center; width: 46px; height: 48px; border: none; background: transparent; cursor: pointer; transition: background 100ms; }
	.wc-btn:hover { background: var(--wash-hover); }
	.wc-btn--close:hover { background: rgba(196,43,28,0.12); }
	.statusbar { display: flex; align-items: center; gap: 16px; height: 30px; padding: 0 16px; border-top: 1px solid var(--border); flex-shrink: 0; font-size: 11.5px; color: var(--ink-muted); }
	.status-item { white-space: nowrap; }
	.status-spacer { flex: 1; }
	.status-item--secure { display: flex; align-items: center; gap: 5px; }
</style>

﻿
