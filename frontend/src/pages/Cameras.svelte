<script lang="ts">
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { GetCameras } from '../../bindings/swail/ajservice';
	import type { MockCamera } from '../../bindings/swail/models';

	const fallback: MockCamera[] = [
		{ id: 'cam-front', label: 'Dash front', status: 'online', scene: 'roadway' },
		{ id: 'cam-rear', label: 'Dash rear', status: 'online', scene: 'roadway' },
		{ id: 'body-3', label: 'Body cam 3', status: 'online', scene: 'patrol' },
		{ id: 'fixed-north', label: 'Fixed north gate', status: 'degraded', scene: 'perimeter' }
	];

	let cameras = $state<MockCamera[]>(fallback);
	let loading = $state(true);
	/** Always keep one camera selected so the player is the primary view */
	let selectedId = $state<string>('cam-front');

	const onlineCount = $derived(cameras.filter((c) => c.status === 'online').length);
	const degradedCount = $derived(cameras.filter((c) => c.status === 'degraded').length);
	const offlineCount = $derived(cameras.filter((c) => c.status === 'offline').length);
	const selected = $derived(cameras.find((c) => c.id === selectedId) ?? cameras[0] ?? null);

	async function load() {
		loading = true;
		try {
			const list = await GetCameras();
			if (list && list.length) {
				cameras = list;
				if (!list.find((c) => c.id === selectedId)) selectedId = list[0].id;
			}
		} catch {
			/* keep fallback */
		} finally {
			loading = false;
		}
	}

	function select(id: string) {
		selectedId = id;
	}

	onMount(() => {
		load();
	});
</script>

<div class="cameras-page">
	<!-- Single-line Tesla-style command bar -->
	<header class="cmd-bar">
		<div class="cmd-left">
			<button class="back-btn touch" onclick={() => push('/')} aria-label="Back">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
					<polyline points="15 18 9 12 15 6"></polyline>
				</svg>
			</button>
			<h1 class="page-title"><span class="title-accent">Cameras</span></h1>
			<span class="page-desc">Sources & live tiles</span>
		</div>

		<div class="cmd-stats" aria-live="polite">
			<span class="stat" data-status="online">
				<span class="stat-dot"></span>
				{onlineCount} Online
			</span>
			<span class="stat" data-status="degraded">
				<span class="stat-dot"></span>
				{degradedCount} Degraded
			</span>
			<span class="stat" data-status="offline">
				<span class="stat-dot"></span>
				{offlineCount} Offline
			</span>
		</div>
	</header>

	{#if loading}
		<p class="empty">Loading sources…</p>
	{:else if cameras.length === 0}
		<p class="empty">No cameras registered.</p>
	{:else}
		<div class="stage">
			<!-- Vertical camera selector rail (always visible) -->
			<nav class="selector" aria-label="Camera selector" role="list">
				{#each cameras as c (c.id)}
					<button
						class="sel-item"
						class:is-active={selectedId === c.id}
						data-status={c.status}
						role="listitem"
						onclick={() => select(c.id)}
						aria-pressed={selectedId === c.id}
						aria-label="{c.label} — {c.status}"
					>
						<div class="sel-preview" data-status={c.status}>
							<div class="preview-mesh" aria-hidden="true"></div>
							<span class="status-pill" data-status={c.status}>{c.status}</span>
						</div>
						<div class="sel-meta">
							<strong class="sel-label">{c.label}</strong>
							<span class="sel-id">{c.id}</span>
						</div>
					</button>
				{/each}
			</nav>

			<!-- 16:9 full-stage player -->
			{#if selected}
				<div class="player-wrap">
					<div class="player" data-status={selected.status}>
						<div class="preview-mesh" aria-hidden="true"></div>

						<!-- Top overlay: specialists + connector -->
						<div class="overlay-top">
							<div class="chip-row">
								<span class="chip">alpr</span>
								<span class="chip">object</span>
								<span class="chip muted">connector-camera</span>
								<span class="chip muted">profile · {selected.scene}</span>
							</div>
						</div>

						<!-- Bottom overlay: identity + metrics -->
						<div class="overlay-bottom">
							<div class="player-oneliner">
								<span class="player-id">{selected.id}</span>
								<span class="player-sep">·</span>
								<span class="player-label">{selected.label}</span>
								<span class="player-sep">·</span>
								<span class="player-scene">{selected.scene}</span>
								<span class="status-pill" data-status={selected.status}>{selected.status}</span>
							</div>
							<div class="player-metrics">
								<span>24 FPS</span>
								<span>18 ms</span>
								<span>3 det</span>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.cameras-page {
		width: 100%;
		min-height: 0;
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	/* ----- Command bar ---------------------------------------------------- */
	.cmd-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem 1.25rem;
		margin-bottom: 0.85rem;
		min-height: 3rem;
		flex-shrink: 0;
	}

	.cmd-left {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		min-width: 0;
	}

	.back-btn.touch {
		margin: 0;
		width: 44px;
		height: 44px;
		padding: 0;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		border-radius: 12px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		color: var(--muted);
		flex-shrink: 0;
		transition: color 0.15s ease, border-color 0.15s ease, background 0.15s ease;
	}

	.back-btn.touch svg {
		width: 22px;
		height: 22px;
	}

	.back-btn.touch:hover {
		color: var(--text);
		border-color: rgba(255, 77, 77, 0.45);
		background: var(--glass-strong);
	}

	.page-title {
		margin: 0;
		font-size: clamp(1.35rem, 2.6vw, 1.7rem);
		line-height: 1.15;
		font-weight: 800;
		letter-spacing: -0.02em;
		white-space: nowrap;
	}

	.page-desc {
		font-size: 13px;
		color: var(--muted);
		font-weight: 500;
		white-space: nowrap;
		opacity: 0.9;
	}

	.cmd-stats {
		display: flex;
		flex-wrap: wrap;
		gap: 0.45rem;
		flex-shrink: 0;
	}

	.stat {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.4rem 0.75rem;
		border-radius: 999px;
		font-size: 12.5px;
		font-weight: 700;
		letter-spacing: 0.02em;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		color: var(--muted);
		--wails-draggable: no-drag;
	}

	.stat-dot {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: currentColor;
	}

	.stat[data-status='online'] {
		color: #4ade80;
		border-color: rgba(74, 222, 128, 0.35);
		background: rgba(34, 197, 94, 0.12);
	}
	.stat[data-status='online'] .stat-dot {
		box-shadow: 0 0 8px #4ade80;
		animation: pulse-dot 2.2s ease-in-out infinite;
	}
	.stat[data-status='degraded'] {
		color: #facc15;
		border-color: rgba(250, 204, 21, 0.4);
		background: rgba(250, 204, 21, 0.12);
	}
	.stat[data-status='degraded'] .stat-dot {
		box-shadow: 0 0 8px #facc15;
	}
	.stat[data-status='offline'] {
		color: #f87171;
		border-color: rgba(248, 113, 113, 0.35);
		background: rgba(239, 68, 68, 0.1);
	}

	@keyframes pulse-dot {
		0%, 100% { opacity: 1; transform: scale(1); }
		50% { opacity: 0.55; transform: scale(0.85); }
	}

	.empty {
		color: var(--muted);
		font-size: 14px;
		margin-top: 0.75rem;
	}

	/* ----- Stage: selector rail + 16:9 player ------------------------------ */
	.stage {
		flex: 1;
		display: flex;
		gap: 0.85rem;
		align-items: stretch;
		width: 100%;
		min-height: 0;
	}

	/* Vertical camera selector */
	.selector {
		flex: 0 0 168px;
		width: 168px;
		display: flex;
		flex-direction: column;
		gap: 0.55rem;
		overflow-y: auto;
		padding-right: 2px;
		--wails-draggable: no-drag;
	}

	.sel-item {
		display: flex;
		flex-direction: column;
		padding: 0;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 12px;
		overflow: hidden;
		cursor: pointer;
		color: inherit;
		font: inherit;
		text-align: left;
		transition: border-color 0.18s ease, box-shadow 0.2s ease, transform 0.15s ease;
	}

	.sel-item:hover {
		border-color: rgba(255, 255, 255, 0.22);
		transform: translateX(2px);
	}

	.sel-item.is-active {
		border-color: rgba(255, 77, 77, 0.55);
		box-shadow: 0 0 0 1px rgba(255, 77, 77, 0.3), 0 6px 20px rgba(255, 45, 114, 0.12);
	}

	.sel-item[data-status='online'].is-active {
		border-color: rgba(74, 222, 128, 0.5);
		box-shadow: 0 0 0 1px rgba(74, 222, 128, 0.25), 0 6px 20px rgba(34, 197, 94, 0.12);
	}

	.sel-preview {
		position: relative;
		height: 72px;
		background: linear-gradient(155deg, #141722 0%, #0a0b12 70%);
		overflow: hidden;
	}

	.sel-preview[data-status='online']::after {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 28px rgba(34, 197, 94, 0.14);
		pointer-events: none;
	}

	.sel-preview[data-status='degraded']::after {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 28px rgba(250, 204, 21, 0.16);
		pointer-events: none;
	}

	.sel-preview .status-pill {
		position: absolute;
		top: 6px;
		right: 6px;
		z-index: 1;
	}

	.preview-mesh {
		position: absolute;
		inset: 0;
		background-image:
			linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
			linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
		background-size: 14px 14px;
		opacity: 0.5;
		pointer-events: none;
	}

	.sel-meta {
		padding: 7px 9px 8px;
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.sel-label {
		font-size: 12.5px;
		font-weight: 700;
		letter-spacing: -0.01em;
		line-height: 1.2;
	}

	.sel-id {
		font-size: 10.5px;
		color: var(--muted);
		font-variant-numeric: tabular-nums;
	}

	.status-pill {
		font-size: 9px;
		font-weight: 800;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		padding: 2px 6px;
		border-radius: 999px;
	}

	.status-pill[data-status='online'] {
		background: rgba(34, 197, 94, 0.28);
		color: #4ade80;
		box-shadow: 0 0 10px rgba(74, 222, 128, 0.3);
	}
	.status-pill[data-status='degraded'] {
		background: rgba(250, 204, 21, 0.28);
		color: #facc15;
		box-shadow: 0 0 10px rgba(250, 204, 21, 0.28);
	}
	.status-pill[data-status='offline'] {
		background: rgba(239, 68, 68, 0.25);
		color: #f87171;
	}

	/* 16:9 player stage */
	.player-wrap {
		flex: 1 1 auto;
		min-width: 0;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.player {
		position: relative;
		width: 100%;
		aspect-ratio: 16 / 9;
		max-height: 100%;
		border-radius: 14px;
		overflow: hidden;
		background: linear-gradient(155deg, #0e1018 0%, #06070f 70%);
		border: 1px solid var(--glass-border);
	}

	.player[data-status='online'] {
		box-shadow: inset 0 0 80px rgba(34, 197, 94, 0.1);
	}

	.player[data-status='degraded'] {
		box-shadow: inset 0 0 80px rgba(250, 204, 21, 0.12);
	}

	.player .preview-mesh {
		background-size: 22px 22px;
		opacity: 0.4;
	}

	/* Overlays */
	.overlay-top {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		z-index: 2;
		padding: 12px 16px;
		background: linear-gradient(rgba(0, 0, 0, 0.55), transparent);
		pointer-events: none;
	}

	.chip-row {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
	}

	.chip {
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.03em;
		padding: 4px 10px;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.14);
		color: var(--text);
		text-transform: lowercase;
	}

	.chip.muted {
		color: var(--muted);
		font-weight: 600;
		background: rgba(255, 255, 255, 0.05);
	}

	.overlay-bottom {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 2;
		padding: 14px 18px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		background: linear-gradient(transparent, rgba(0, 0, 0, 0.72));
		pointer-events: none;
	}

	.player-oneliner {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-wrap: wrap;
		min-width: 0;
	}

	.player-id {
		font-size: 13px;
		font-weight: 700;
		letter-spacing: 0.04em;
		font-variant-numeric: tabular-nums;
		color: rgba(244, 246, 251, 0.85);
	}

	.player-sep {
		color: var(--muted);
		opacity: 0.6;
	}

	.player-label {
		font-size: 14px;
		font-weight: 700;
	}

	.player-scene {
		font-size: 13px;
		color: var(--muted);
		text-transform: capitalize;
	}

	.player-metrics {
		display: flex;
		gap: 14px;
		font-size: 12.5px;
		font-weight: 600;
		color: var(--muted);
		font-variant-numeric: tabular-nums;
		flex-shrink: 0;
	}

	/* ----- Narrow / mobile ------------------------------------------------ */
	@media (max-width: 720px) {
		.stage {
			flex-direction: column;
		}

		.selector {
			flex: none;
			width: 100%;
			flex-direction: row;
			overflow-x: auto;
			overflow-y: hidden;
			gap: 0.5rem;
			padding-bottom: 4px;
		}

		.sel-item {
			flex: 0 0 140px;
		}

		.sel-preview {
			height: 64px;
		}

		.page-desc {
			display: none;
		}

		.player {
			border-radius: 12px;
		}
	}
</style>
