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
	let selectedId = $state<string | null>(null);

	const onlineCount = $derived(cameras.filter((c) => c.status === 'online').length);
	const degradedCount = $derived(cameras.filter((c) => c.status === 'degraded').length);
	const offlineCount = $derived(cameras.filter((c) => c.status === 'offline').length);

	async function load() {
		loading = true;
		try {
			const list = await GetCameras();
			if (list && list.length) cameras = list;
		} catch {
			/* keep fallback */
		} finally {
			loading = false;
		}
	}

	function select(id: string) {
		selectedId = selectedId === id ? null : id;
	}

	onMount(load);
</script>

<div class="cameras-page">
	<button class="back-btn" onclick={() => push('/')}>← Back</button>

	<!-- Command bar -->
	<header class="cmd-bar">
		<div class="cmd-left">
			<h1 class="title"><span class="title-accent">Cameras</span></h1>
			<p class="subtitle">Sources registered with scene-gate / connector-camera</p>
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
		<div class="cam-grid" role="list">
			{#each cameras as c (c.id)}
				<button
					class="cam-tile"
					class:is-selected={selectedId === c.id}
					data-status={c.status}
					role="listitem"
					onclick={() => select(c.id)}
					aria-pressed={selectedId === c.id}
					aria-label="{c.label} — {c.status}"
				>
					<div class="cam-preview" data-status={c.status}>
						<div class="preview-overlay">
							<span class="cam-id">{c.id}</span>
							<span class="status-pill" data-status={c.status}>{c.status}</span>
						</div>
						<!-- subtle scan / mesh placeholder until live frames arrive -->
						<div class="preview-mesh" aria-hidden="true"></div>
					</div>

					<div class="cam-meta">
						<strong class="cam-label">{c.label}</strong>
						<span class="scene">{c.scene}</span>
					</div>
				</button>
			{/each}
		</div>
	{/if}
</div>

<style>
	/* Full-bleed relative to the narrow App container */
	.cameras-page {
		width: 100vw;
		max-width: none;
		margin-left: calc(-50vw + 50%);
		padding: 0 max(1.5rem, env(safe-area-inset-right)) 2rem max(1.5rem, env(safe-area-inset-left));
		box-sizing: border-box;
	}

	/* ----- Command bar ---------------------------------------------------- */
	.cmd-bar {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		justify-content: space-between;
		gap: 1rem 2rem;
		margin-bottom: 1.75rem;
	}

	.cmd-left .title {
		margin-bottom: 0.35rem;
	}

	.cmd-left .subtitle {
		margin: 0;
		font-size: 0.95rem;
	}

	.cmd-stats {
		display: flex;
		flex-wrap: wrap;
		gap: 0.6rem;
	}

	.stat {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.35rem 0.75rem;
		border-radius: 999px;
		font-size: 12px;
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
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.55;
			transform: scale(0.85);
		}
	}

	/* ----- Grid ------------------------------------------------------------ */
	.cam-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
		gap: 1rem;
		width: 100%;
	}

	/* ----- Card ------------------------------------------------------------ */
	.cam-tile {
		display: flex;
		flex-direction: column;
		text-align: left;
		padding: 0;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 14px;
		overflow: hidden;
		cursor: pointer;
		color: inherit;
		font: inherit;
		transition:
			transform 0.18s ease,
			box-shadow 0.22s ease,
			border-color 0.2s ease;
		--wails-draggable: no-drag;
	}

	.cam-tile:hover {
		transform: translateY(-3px);
		border-color: rgba(255, 255, 255, 0.22);
	}

	.cam-tile.is-selected {
		border-color: rgba(255, 77, 77, 0.55);
		box-shadow: 0 0 0 1px rgba(255, 77, 77, 0.35);
	}

	/* Status glow on the whole card */
	.cam-tile[data-status='online'] {
		box-shadow: 0 0 0 1px rgba(74, 222, 128, 0.12), 0 8px 28px rgba(34, 197, 94, 0.08);
	}
	.cam-tile[data-status='online']:hover {
		box-shadow: 0 0 0 1px rgba(74, 222, 128, 0.35), 0 12px 36px rgba(34, 197, 94, 0.18);
	}
	.cam-tile[data-status='degraded'] {
		box-shadow: 0 0 0 1px rgba(250, 204, 21, 0.18), 0 8px 28px rgba(250, 204, 21, 0.08);
	}
	.cam-tile[data-status='degraded']:hover {
		box-shadow: 0 0 0 1px rgba(250, 204, 21, 0.4), 0 12px 36px rgba(250, 204, 21, 0.16);
	}
	.cam-tile[data-status='offline'] {
		opacity: 0.72;
		filter: saturate(0.65);
	}

	/* Preview area */
	.cam-preview {
		position: relative;
		height: 140px;
		background: linear-gradient(155deg, #141722 0%, #0a0b12 70%);
		overflow: hidden;
	}

	.preview-mesh {
		position: absolute;
		inset: 0;
		background-image:
			linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
			linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
		background-size: 18px 18px;
		opacity: 0.55;
		pointer-events: none;
	}

	.cam-preview[data-status='online']::after {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 40px rgba(34, 197, 94, 0.12);
		pointer-events: none;
		animation: glow-breathe 3.5s ease-in-out infinite;
	}

	.cam-preview[data-status='degraded']::after {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 36px rgba(250, 204, 21, 0.14);
		pointer-events: none;
	}

	@keyframes glow-breathe {
		0%,
		100% {
			opacity: 0.7;
		}
		50% {
			opacity: 1;
		}
	}

	.preview-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		padding: 10px 12px;
		z-index: 1;
		background: linear-gradient(transparent 40%, rgba(0, 0, 0, 0.55));
	}

	.cam-id {
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.05em;
		color: rgba(244, 246, 251, 0.75);
		font-variant-numeric: tabular-nums;
	}

	.status-pill {
		font-size: 10px;
		font-weight: 800;
		text-transform: uppercase;
		letter-spacing: 0.07em;
		padding: 3px 8px;
		border-radius: 999px;
	}

	.status-pill[data-status='online'] {
		background: rgba(34, 197, 94, 0.28);
		color: #4ade80;
		box-shadow: 0 0 12px rgba(74, 222, 128, 0.35);
	}
	.status-pill[data-status='degraded'] {
		background: rgba(250, 204, 21, 0.28);
		color: #facc15;
		box-shadow: 0 0 12px rgba(250, 204, 21, 0.3);
	}
	.status-pill[data-status='offline'] {
		background: rgba(239, 68, 68, 0.25);
		color: #f87171;
	}

	/* Meta */
	.cam-meta {
		padding: 12px 14px 14px;
		display: flex;
		flex-direction: column;
		gap: 3px;
	}

	.cam-label {
		font-size: 15px;
		font-weight: 700;
		letter-spacing: -0.01em;
	}

	.scene {
		font-size: 12.5px;
		color: var(--muted);
		text-transform: capitalize;
	}

	.empty {
		color: var(--muted);
		font-size: 14px;
		margin-top: 1rem;
	}

	/* Mobile */
	@media (max-width: 640px) {
		.cameras-page {
			margin-left: 0;
			width: 100%;
			padding-left: 0;
			padding-right: 0;
		}

		.cam-grid {
			grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
			gap: 0.75rem;
		}

		.cam-preview {
			height: 110px;
		}

		.cmd-bar {
			flex-direction: column;
			align-items: flex-start;
		}
	}
</style>
