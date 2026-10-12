<script lang="ts">
	import { push } from 'svelte-spa-router';

	// svelte-spa-router passes params via props
	let { params }: { params?: { id?: string } } = $props();
	const id = $derived(params?.id ?? 'unknown');

	// PR-5: fetch presigned URL from object-storage GET /v1/objects/{id}/presign
	const mockMeta = {
		kind: 'camera_recording',
		source: 'cam-front',
		contentType: 'video/mp4',
		size: '2.4 MB',
		ts: '2026-10-12 14:22:01 UTC'
	};
</script>

<button class="back-btn" onclick={() => push('/')}>← Back</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">Evidence</span></h1>
	<p class="subtitle">object_id <code>{id}</code> · will stream via presigned URL (PR-5).</p>
</div>

<div class="player-shell">
	<div class="player-placeholder" role="img" aria-label="Video placeholder">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
			<polygon points="5 3 19 12 5 21 5 3"></polygon>
		</svg>
		<span>Presigned media will appear here</span>
	</div>

	<dl class="meta">
		<div><dt>Kind</dt><dd>{mockMeta.kind}</dd></div>
		<div><dt>Source</dt><dd>{mockMeta.source}</dd></div>
		<div><dt>Type</dt><dd>{mockMeta.contentType}</dd></div>
		<div><dt>Size</dt><dd>{mockMeta.size}</dd></div>
		<div><dt>Captured</dt><dd>{mockMeta.ts}</dd></div>
	</dl>

	<p class="note">
		Requires object-storage <code>GET /v1/objects/&#123;id&#125;/presign</code> + CORS (core PR-3).
	</p>
</div>

<style>
	.player-shell {
		width: 100%;
		max-width: 640px;
	}
	.player-placeholder {
		aspect-ratio: 16 / 9;
		background: #0a0b12;
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;
		color: var(--muted);
		font-size: 13px;
	}
	.meta {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
		gap: 10px;
		margin: var(--s-3) 0;
	}
	.meta div {
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 8px;
		padding: 8px 10px;
	}
	.meta dt {
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--muted);
		margin: 0 0 2px;
	}
	.meta dd {
		margin: 0;
		font-size: 13.5px;
		font-weight: 600;
	}
	.note {
		font-size: 12.5px;
		color: var(--muted);
		line-height: 1.5;
	}
	code {
		font-size: 0.92em;
		background: rgba(255, 255, 255, 0.06);
		padding: 1px 5px;
		border-radius: 4px;
	}
</style>
