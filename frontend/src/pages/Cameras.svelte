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
	const selected = $derived(cameras.find((c) => c.id === selectedId) ?? null);

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

	function closeInspector() {
		selectedId = null;
	}

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && selectedId) {
			e.preventDefault();
			closeInspector();
		}
	}

	onMount(() => {
		load();
		window.addEventListener('keydown', onKeydown);
		return () => window.removeEventListener('keydown', onKeydown);
	});
</script>

<div class="cameras-page" class:has-inspector={!!selected}>
	<!-- Single-line Tesla-style command bar: back · title · description · stats -->
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
			<!-- Left-sliding inspector sidebar -->
			{#if selected}
				<aside class="inspector" aria-label="Camera details">
					<div class="insp-head">
						<div class="insp-title-row">
							<h2 class="insp-title">{selected.label}</h2>
							<span class="status-pill" data-status={selected.status}>{selected.status}</span>
						</div>
						<button class="insp-close touch" onclick={closeInspector} aria-label="Close inspector">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
								<line x1="18" y1="6" x2="6" y2="18"></line>
								<line x1="6" y1="6" x2="18" y2="18"></line>
							</svg>
						</button>
					</div>

					<div class="insp-preview" data-status={selected.status}>
						<div class="preview-mesh" aria-hidden="true"></div>
						<span class="insp-preview-label">{selected.id}</span>
					</div>

					<dl class="insp-meta">
						<div class="meta-row">
							<dt>Scene</dt>
							<dd class="cap">{selected.scene}</dd>
						</div>
						<div class="meta-row">
							<dt>Source ID</dt>
							<dd class="mono">{selected.id}</dd>
						</div>
						<div class="meta-row">
							<dt>Connector</dt>
							<dd>connector-camera</dd>
						</div>
						<div class="meta-row">
							<dt>Scene-gate</dt>
							<dd>profile · {selected.scene}</dd>
						</div>
						<div class="meta-row">
							<dt>Specialists</dt>
							<dd>alpr · object</dd>
						</div>
					</dl>

					<div class="insp-metrics">
						<div class="metric">
							<span class="metric-val">24</span>
							<span class="metric-lbl">FPS</span>
						</div>
						<div class="metric">
							<span class="metric-val">18ms</span>
							<span class="metric-lbl">Latency</span>
						</div>
						<div class="metric">
							<span class="metric-val">3</span>
							<span class="metric-lbl">Detections</span>
						</div>
					</div>

					<div class="insp-actions">
						<button class="act primary touch" type="button">View live</button>
						<button class="act touch" type="button">Pin</button>
						<button class="act touch" type="button" onclick={closeInspector}>Close</button>
					</div>
				</aside>
			{/if}

			<div class="main-stage">
				{#if selected}
					<!-- Full-screen player with overlay one-liner -->
					<div class="player" data-status={selected.status}>
						<div class="preview-mesh" aria-hidden="true"></div>
						<div class="player-overlay">
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
						<button class="player-close touch" onclick={closeInspector} aria-label="Close player">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
								<line x1="18" y1="6" x2="6" y2="18"></line>
								<line x1="6" y1="6" x2="18" y2="18"></line>
							</svg>
						</button>
					</div>
				{:else}
					<div class="cam-grid" role="list">
						{#each cameras as c (c.id)}
							<button
								class="cam-tile"
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

	.stage {
		flex: 1;
		display: flex;
		gap: 0;
		align-items: stretch;
		width: 100%;
		min-height: 0;
		position: relative;
	}

	.main-stage {
		flex: 1 1 auto;
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	.cam-grid {
		flex: 1;
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 0.9rem;
		align-content: start;
	}

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
		transition: transform 0.18s ease, box-shadow 0.22s ease, border-color 0.2s ease;
		--wails-draggable: no-drag;
	}

	.cam-tile:hover {
		transform: translateY(-3px);
		border-color: rgba(255, 255, 255, 0.22);
	}

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

	.cam-preview {
		position: relative;
		height: 130px;
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
		0%, 100% { opacity: 0.7; }
		50% { opacity: 1; }
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

	.cam-meta {
		padding: 10px 12px 12px;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.cam-label {
		font-size: 14px;
		font-weight: 700;
		letter-spacing: -0.01em;
	}

	.scene {
		font-size: 12px;
		color: var(--muted);
		text-transform: capitalize;
	}

	.empty {
		color: var(--muted);
		font-size: 14px;
		margin-top: 0.75rem;
	}

	.inspector {
		flex: 0 0 300px;
		width: 300px;
		max-width: 36vw;
		background: rgba(8, 10, 18, 0.82);
		border-right: 1px solid var(--glass-border);
		border-radius: 0 14px 14px 0;
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 14px;
		-webkit-backdrop-filter: blur(20px);
		backdrop-filter: blur(20px);
		animation: sidebar-in 0.28s cubic-bezier(0.22, 0.68, 0.2, 1);
		--wails-draggable: no-drag;
		overflow-y: auto;
	}

	@keyframes sidebar-in {
		from { opacity: 0; transform: translateX(-100%); }
		to { opacity: 1; transform: translateX(0); }
	}

	.insp-head {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 8px;
	}

	.insp-title-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 8px;
		min-width: 0;
	}

	.insp-title {
		margin: 0;
		font-size: 1.1rem;
		font-weight: 800;
		letter-spacing: -0.01em;
	}

	.insp-close.touch {
		flex-shrink: 0;
		width: 40px;
		height: 40px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		border: 1px solid var(--glass-border);
		border-radius: 10px;
		background: transparent;
		color: var(--muted);
		cursor: pointer;
		transition: color 0.15s ease, border-color 0.15s ease;
	}

	.insp-close.touch svg {
		width: 18px;
		height: 18px;
	}

	.insp-close.touch:hover {
		color: var(--text);
		border-color: rgba(255, 77, 77, 0.45);
	}

	.insp-preview {
		position: relative;
		height: 140px;
		border-radius: 10px;
		overflow: hidden;
		background: linear-gradient(155deg, #141722 0%, #0a0b12 70%);
	}

	.insp-preview[data-status='online'] {
		box-shadow: inset 0 0 48px rgba(34, 197, 94, 0.14);
	}

	.insp-preview[data-status='degraded'] {
		box-shadow: inset 0 0 48px rgba(250, 204, 21, 0.16);
	}

	.insp-preview-label {
		position: absolute;
		left: 10px;
		bottom: 10px;
		z-index: 1;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.05em;
		color: rgba(244, 246, 251, 0.8);
	}

	.insp-meta {
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 7px;
	}

	.meta-row {
		display: flex;
		justify-content: space-between;
		gap: 12px;
		font-size: 13px;
	}

	.meta-row dt { color: var(--muted); margin: 0; }
	.meta-row dd { margin: 0; font-weight: 600; text-align: right; }
	.meta-row dd.cap { text-transform: capitalize; }
	.meta-row dd.mono {
		font-variant-numeric: tabular-nums;
		font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
		font-size: 12.5px;
	}

	.insp-metrics {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 8px;
	}

	.metric {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2px;
		padding: 10px 4px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 10px;
	}

	.metric-val {
		font-size: 16px;
		font-weight: 800;
		font-variant-numeric: tabular-nums;
	}

	.metric-lbl {
		font-size: 10px;
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: var(--muted);
	}

	.insp-actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
		margin-top: auto;
	}

	.act.touch {
		min-height: 44px;
		padding: 10px 14px;
		border-radius: 11px;
		border: 1px solid var(--glass-border);
		background: var(--glass);
		color: var(--text);
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: border-color 0.15s ease, background 0.15s ease;
		--wails-draggable: no-drag;
	}

	.act.touch:hover {
		border-color: rgba(255, 255, 255, 0.28);
		background: var(--glass-strong);
	}

	.act.primary {
		background: var(--grad);
		border-color: transparent;
		color: #fff;
		box-shadow: 0 4px 14px rgba(255, 45, 114, 0.28);
	}

	.act.primary:hover {
		box-shadow: 0 6px 18px rgba(255, 45, 114, 0.4);
	}

	.player {
		position: relative;
		flex: 1;
		min-height: 320px;
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

	.player-overlay {
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

	.player-sep { color: var(--muted); opacity: 0.6; }
	.player-label { font-size: 14px; font-weight: 700; }
	.player-scene { font-size: 13px; color: var(--muted); text-transform: capitalize; }

	.player-metrics {
		display: flex;
		gap: 12px;
		font-size: 12px;
		font-weight: 600;
		color: var(--muted);
		font-variant-numeric: tabular-nums;
	}

	.player-close.touch {
		position: absolute;
		top: 12px;
		right: 12px;
		z-index: 3;
		width: 44px;
		height: 44px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		border: 1px solid var(--glass-border);
		border-radius: 12px;
		background: rgba(10, 12, 20, 0.65);
		color: var(--muted);
		cursor: pointer;
		-webkit-backdrop-filter: blur(8px);
		backdrop-filter: blur(8px);
		transition: color 0.15s ease, border-color 0.15s ease;
		--wails-draggable: no-drag;
	}

	.player-close.touch svg { width: 20px; height: 20px; }

	.player-close.touch:hover {
		color: var(--text);
		border-color: rgba(255, 77, 77, 0.5);
	}

	@media (max-width: 720px) {
		.stage { flex-direction: column; }
		.inspector {
			flex: none;
			width: 100%;
			max-width: none;
			border-right: none;
			border-bottom: 1px solid var(--glass-border);
			border-radius: 14px 14px 0 0;
			animation: none;
		}
		.cam-grid { grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); }
		.cam-preview { height: 110px; }
		.page-desc { display: none; }
		.player { min-height: 240px; }
	}
</style>
