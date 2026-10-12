<script lang="ts">
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import { newDocumentTrigger } from '#lib/stores';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Opaque backdrop — overrides Wails WebviewIsTransparent -->
<div class="backdrop">
	<div class="shell">
		<aside class="sidebar">
			<div class="sidebar-header">
				<span class="app-icon">🖊️</span>
				<span class="app-name">ink-mark</span>
			</div>

			<nav class="nav">
				<button class="nav-item nav-item--active" type="button"
					onclick={() => newDocumentTrigger.update(n => n + 1)}>
					<i class="ph-duotone ph-note-pencil nav-icon"></i>
					<span>New Document</span>
				</button>
				<button class="nav-item" type="button">
					<i class="ph-duotone ph-folder-open nav-icon"></i>
					<span>All Documents</span>
				</button>
			</nav>

			<div class="sidebar-spacer"></div>

			<nav class="nav">
				<button class="nav-item" type="button">
					<i class="ph-duotone ph-gear-six nav-icon"></i>
					<span>Settings</span>
				</button>
			</nav>
		</aside>

		<main class="content">
			{@render children()}
		</main>
	</div>
</div>

<style>
	/* Opaque full-bleed backdrop */
	.backdrop {
		position: fixed;
		inset: 0;
		background: #F3F3F3;
	}

	.shell {
		position: relative;
		display: grid;
		grid-template-columns: 240px 1fr;
		height: 100vh;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		color: #1A1A1A;
	}

	/* ── Sidebar ─────────────────────────────────────────── */
	.sidebar {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: rgba(243, 243, 243, 0.96);
		border-right: 1px solid #E5E5E5;
		padding: 12px 8px;
	}

	.sidebar-header {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 12px 20px;
	}

	.app-icon { font-size: 18px; line-height: 1; }

	.app-name {
		font-size: 14px;
		font-weight: 600;
		color: #1A1A1A;
		letter-spacing: -0.01em;
	}

	.nav { display: flex; flex-direction: column; gap: 1px; }
	.sidebar-spacer { flex: 1; }

	.nav-item {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 7px 12px;
		border: none;
		border-radius: 6px;
		background: transparent;
		color: #606060;
		font-family: 'Segoe UI Variable', 'Segoe UI', system-ui, sans-serif;
		font-size: 13px;
		font-weight: 400;
		cursor: pointer;
		text-align: left;
		outline: none;
		position: relative;
		transition: background 120ms ease, color 120ms ease;
	}

	.nav-item:hover { background: #F0F0F0; color: #1A1A1A; }
	.nav-item:active { background: #E7E7E7; transform: scale(0.99); }
	.nav-item:focus-visible { box-shadow: 0 0 0 2px #005FB8; }

	.nav-item--active {
		background: rgba(15, 108, 189, 0.10);
		color: #0F6CBD;
		font-weight: 600;
	}

	.nav-item--active::before {
		content: '';
		position: absolute;
		left: 0;
		top: 50%;
		transform: translateY(-50%);
		width: 3px;
		height: 16px;
		background: #0F6CBD;
		border-radius: 0 4px 4px 0;
	}

	.nav-icon { font-size: 16px; flex-shrink: 0; color: inherit; }

	/* ── Content ─────────────────────────────────────────── */
	.content {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: #FFFFFF;
		overflow: hidden;
	}
</style>


