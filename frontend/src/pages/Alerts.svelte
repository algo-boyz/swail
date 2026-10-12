<script lang="ts">
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { GetRecentAlerts } from '../../bindings/swail/ajservice';
	import type { MockAlert } from '../../bindings/swail/ajservice';
	import { _ } from '../lib/i18n';

	let alerts = $state<MockAlert[]>([]);
	let loading = $state(true);
	let error = $state('');

	async function load() {
		loading = true;
		error = '';
		try {
			alerts = (await GetRecentAlerts()) ?? [];
		} catch (e: any) {
			error = e?.message ?? String(e);
			alerts = [];
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

<button class="back-btn" onclick={() => push('/')}>{$_('common.back')}</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">{$_('alerts.title')}</span></h1>
	<p class="subtitle">
		{$_('alerts.subtitle')}
		<button class="linkish inline" onclick={load} disabled={loading}>{$_('common.refresh')}</button>
	</p>
</div>

{#if error}
	<p class="err">{error}</p>
{/if}

<div class="alert-list">
	{#if loading}
		<p class="empty">{$_('common.loading')}</p>
	{:else}
		{#each alerts as a}
			<div class="alert-card" data-kind={a.kind}>
				<div class="alert-head">
					<span class="kind-badge">{a.kind}</span>
					<span class="ts">{a.timestamp}</span>
				</div>
				<p class="msg">{a.message}</p>
				<div class="alert-foot">
					<span class="cam">{a.camera_id}</span>
					<span class="conf">{Math.round(a.confidence * 100)}%</span>
					{#if a.object_id}
						<button class="linkish" onclick={() => push(`/evidence/${a.object_id}`)}
							>{$_('alerts.evidence')}</button
						>
					{/if}
				</div>
			</div>
		{:else}
			<p class="empty">{$_('alerts.empty')}</p>
		{/each}
	{/if}
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
	.alert-card[data-kind='plate'],
	.alert-card[data-kind='detection'] {
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
	.linkish.inline {
		margin-left: 10px;
		font-size: 13px;
	}
	.empty {
		color: var(--muted);
		font-size: 14px;
	}
	.err {
		color: #f87171;
		font-size: 13px;
	}
</style>
