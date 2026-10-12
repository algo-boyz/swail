<script lang="ts">
	import { push } from 'svelte-spa-router';

	const alerts = [
		{
			id: 'a-001',
			kind: 'plate',
			camera: 'cam-front',
			message: '7ABC123 · Mercedes C-Class',
			confidence: 0.94,
			ts: '14:22:01',
			objectId: 'obj-rec-001'
		},
		{
			id: 'a-002',
			kind: 'scene',
			camera: 'cam-rear',
			message: 'roadway / driving',
			confidence: 0.88,
			ts: '14:21:44'
		},
		{
			id: 'a-003',
			kind: 'officer-down',
			camera: 'body-3',
			message: 'Possible officer-down (audio intent)',
			confidence: 0.91,
			ts: '14:18:12',
			objectId: 'obj-aud-003'
		}
	];
</script>

<button class="back-btn" onclick={() => push('/')}>← Back</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">Alerts</span></h1>
	<p class="subtitle">Recent detections & audio intents (mock · PR-2 will poll voice-query / NATS).</p>
</div>

<div class="alert-list">
	{#each alerts as a}
		<div class="alert-card" data-kind={a.kind}>
			<div class="alert-head">
				<span class="kind-badge">{a.kind}</span>
				<span class="ts">{a.ts}</span>
			</div>
			<p class="msg">{a.message}</p>
			<div class="alert-foot">
				<span class="cam">{a.camera}</span>
				<span class="conf">{Math.round(a.confidence * 100)}%</span>
				{#if a.objectId}
					<button class="linkish" onclick={() => push(`/evidence/${a.objectId}`)}
						>Evidence →</button
					>
				{/if}
			</div>
		</div>
	{/each}
</div>

<style>
	.alert-list {
		display: flex;
		flex-direction: column;
		gap: 10px;
		width: 100%;
		max-width: 560px;
	}
	.alert-card {
		padding: 12px 14px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
		border-left-width: 3px;
	}
	.alert-card[data-kind='officer-down'] {
		border-left-color: #ef4444;
	}
	.alert-card[data-kind='plate'] {
		border-left-color: #3b82f6;
	}
	.alert-card[data-kind='scene'] {
		border-left-color: #a78bfa;
	}
	.alert-head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 6px;
	}
	.kind-badge {
		font-size: 10px;
		font-weight: 800;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: #ff7a8f;
	}
	.ts {
		font-size: 12px;
		color: var(--muted);
		font-variant-numeric: tabular-nums;
	}
	.msg {
		margin: 0;
		font-size: 14.5px;
		font-weight: 600;
	}
	.alert-foot {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-top: 8px;
		font-size: 12.5px;
		color: var(--muted);
	}
	.conf {
		color: #4ade80;
		font-weight: 700;
	}
	.linkish {
		margin-left: auto;
		background: none;
		border: none;
		color: #ff7a8f;
		font-size: 12.5px;
		font-weight: 600;
		cursor: pointer;
		padding: 0;
		--wails-draggable: no-drag;
	}
</style>
