<script lang="ts">
	import { push } from 'svelte-spa-router';

	const cameras = [
		{ id: 'cam-front', label: 'Dash front', status: 'online', scene: 'roadway' },
		{ id: 'cam-rear', label: 'Dash rear', status: 'online', scene: 'roadway' },
		{ id: 'body-3', label: 'Body cam 3', status: 'online', scene: 'patrol' },
		{ id: 'fixed-north', label: 'Fixed north gate', status: 'degraded', scene: 'perimeter' }
	];
</script>

<button class="back-btn" onclick={() => push('/')}>← Back</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">Cameras</span></h1>
	<p class="subtitle">Sources registered with scene-gate / connector-camera (mock).</p>
</div>

<div class="cam-grid">
	{#each cameras as c}
		<div class="cam-tile">
			<div class="cam-preview" data-status={c.status}>
				<span class="cam-id">{c.id}</span>
			</div>
			<div class="cam-meta">
				<strong>{c.label}</strong>
				<span class="status-pill" data-status={c.status}>{c.status}</span>
				<span class="scene">{c.scene}</span>
			</div>
		</div>
	{/each}
</div>

<style>
	.cam-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
		gap: 12px;
		width: 100%;
		max-width: 720px;
	}
	.cam-tile {
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
	overflow: hidden;
	}
	.cam-preview {
		height: 100px;
		background: linear-gradient(145deg, #12141f, #0a0b12);
		display: flex;
		align-items: flex-end;
		padding: 8px;
	}
	.cam-id {
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.04em;
		color: var(--muted);
	}
	.cam-meta {
		padding: 10px 12px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.cam-meta strong {
		font-size: 14px;
	}
	.status-pill {
		align-self: flex-start;
		font-size: 10px;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		padding: 2px 7px;
		border-radius: 999px;
	}
	.status-pill[data-status='online'] {
		background: rgba(34, 197, 94, 0.2);
		color: #4ade80;
	}
	.status-pill[data-status='degraded'] {
		background: rgba(250, 204, 21, 0.2);
		color: #facc15;
	}
	.status-pill[data-status='offline'] {
		background: rgba(239, 68, 68, 0.2);
		color: #f87171;
	}
	.scene {
		font-size: 12px;
		color: var(--muted);
	}
</style>
