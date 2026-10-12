<script lang="ts">
	import { push } from 'svelte-spa-router';
	import { Query as RunQuery, Health } from '../../bindings/swail/ajservice';
	import type { QueryResultRow } from '../../bindings/swail/ajservice';

	const tools = [
		{ id: 'plates', label: 'Plates', hint: 'query_recent_plates · make_model, time_window' },
		{ id: 'objects', label: 'Objects', hint: 'query_objects · class, track_id, camera' },
		{ id: 'scenes', label: 'Scenes', hint: 'query_scenes · L1/L2, camera' },
		{ id: 'transcripts', label: 'Transcripts', hint: 'query_transcripts · contains, window' },
		{ id: 'detections', label: 'Detections', hint: 'query_detections · type, camera' }
	];

	let active = $state('plates');
	let windowSec = $state(30);
	let limit = $state(5);
	let filter = $state('');
	let results = $state<QueryResultRow[]>([]);
	let loading = $state(false);
	let error = $state('');
	let ran = $state(false);
	let vqStatus = $state('checking…');

	async function checkHealth() {
		try {
			const h = await Health();
			const vq = h?.voice_query as Record<string, unknown> | undefined;
			if (vq?.ok) {
				vqStatus = vq.demo_mode ? 'voice-query · demo' : 'voice-query · live';
			} else {
				vqStatus = 'voice-query offline';
			}
		} catch {
			vqStatus = 'bridge unavailable';
		}
	}

	async function runQuery() {
		loading = true;
		error = '';
		ran = true;
		try {
			const rows = await RunQuery({
				tool: active,
				filter: filter,
				time_window_sec: windowSec,
				limit: limit
			});
			results = rows ?? [];
		} catch (e: any) {
			error = e?.message ?? String(e);
			results = [];
		} finally {
			loading = false;
			checkHealth();
		}
	}

	checkHealth();
</script>

<button class="back-btn" onclick={() => push('/')}>← Back</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">Query</span></h1>
	<p class="subtitle">
		Proxied to voice-query-service <code>/v1/query/*</code>
		<span class="status-chip">{vqStatus}</span>
	</p>
</div>

<div class="query-layout">
	<div class="tool-tabs">
		{#each tools as t}
			<button
				class="pill"
				class:active={active === t.id}
				onclick={() => {
					active = t.id;
					results = [];
					error = '';
					ran = false;
				}}
			>
				{t.label}
			</button>
		{/each}
	</div>
	<p class="tool-hint">{tools.find((t) => t.id === active)?.hint}</p>

	<div class="form-row">
		<label class="field">
			<span class="control-label">Filter</span>
			<input class="input" type="text" bind:value={filter} placeholder="make_model / class / phrase…" />
		</label>
		<label class="field narrow">
			<span class="control-label">Window (s)</span>
			<input class="input" type="number" bind:value={windowSec} min="5" max="3600" />
		</label>
		<label class="field narrow">
			<span class="control-label">Limit</span>
			<input class="input" type="number" bind:value={limit} min="1" max="50" />
		</label>
		<button class="btn" onclick={runQuery} disabled={loading}>
			{loading ? '…' : 'Run'}
		</button>
	</div>

	{#if error}
		<p class="err">{error}</p>
	{/if}

	<div class="results">
		{#if loading}
			<p class="empty">Loading…</p>
		{:else}
			{#each results as r}
				<div class="result-card">
					<div class="result-top">
						<strong>{r.title}</strong>
						{#if r.confidence != null}
							<span class="conf">{Math.round(r.confidence * 100)}%</span>
						{/if}
					</div>
					<p class="result-meta">{r.meta}</p>
					{#if r.object_id}
						<button class="linkish" onclick={() => push(`/evidence/${r.object_id}`)}
							>Play evidence →</button
						>
					{/if}
				</div>
			{:else}
				<p class="empty">
					{#if ran && !error}
						No hits in this window.
					{:else if !error}
						Run a query against voice-query-service (DEMO_MODE=true works offline).
					{/if}
				</p>
			{/each}
		{/if}
	</div>
</div>

<style>
	.query-layout {
		width: 100%;
		max-width: 640px;
	}
	.status-chip {
		margin-left: 8px;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: #ff7a8f;
		background: rgba(255, 45, 114, 0.14);
		border: 1px solid rgba(255, 45, 114, 0.3);
		border-radius: 999px;
		padding: 2px 8px;
	}
	.tool-tabs {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
		margin-bottom: 8px;
	}
	.pill {
		padding: 5px 12px;
		font-size: 13px;
		font-weight: 600;
		border-radius: 999px;
		border: 1px solid var(--glass-border);
		background: var(--glass);
		color: var(--muted);
		cursor: pointer;
		--wails-draggable: no-drag;
	}
	.pill.active {
		background: rgba(255, 45, 114, 0.18);
		border-color: rgba(255, 45, 114, 0.55);
		color: #ff7a8f;
	}
	.tool-hint {
		font-size: 12px;
		color: var(--muted);
		margin: 0 0 var(--s-3);
	}
	.form-row {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
		align-items: flex-end;
		margin-bottom: var(--s-4);
	}
	.field {
		flex: 1 1 140px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.field.narrow {
		flex: 0 0 90px;
	}
	.control-label {
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--muted);
	}
	.input {
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 8px;
		color: var(--text);
		padding: 8px 10px;
		font-size: 14px;
		outline: none;
		--wails-draggable: no-drag;
	}
	.input:focus {
		border-color: rgba(255, 77, 77, 0.55);
	}
	.btn:disabled {
		opacity: 0.6;
		cursor: wait;
	}
	.err {
		color: #f87171;
		font-size: 13px;
		margin: 0 0 12px;
	}
	.results {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}
	.result-card {
		padding: 12px 14px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
	}
	.result-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.conf {
		font-size: 12px;
		font-weight: 700;
		color: #4ade80;
	}
	.result-meta {
		margin: 4px 0 0;
		font-size: 13px;
		color: var(--muted);
	}
	.linkish {
		margin-top: 8px;
		background: none;
		border: none;
		color: #ff7a8f;
		font-size: 13px;
		font-weight: 600;
		cursor: pointer;
		padding: 0;
		--wails-draggable: no-drag;
	}
	.empty {
		color: var(--muted);
		font-size: 14px;
	}
	code {
		font-size: 0.9em;
		background: rgba(255, 255, 255, 0.06);
		padding: 1px 5px;
		border-radius: 4px;
	}
</style>
