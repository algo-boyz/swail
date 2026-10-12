<script lang="ts">
	import { push } from 'svelte-spa-router';
	import { GetEvidence } from '../../bindings/swail/ajservice';
	import type { EvidenceMeta } from '../../bindings/swail/models';
	import { _ } from '../lib/i18n';

	let { params }: { params?: { id?: string } } = $props();
	const id = $derived(params?.id ?? '');

	let meta = $state<EvidenceMeta | null>(null);
	let loading = $state(true);
	let error = $state('');

	async function load(objectId: string) {
		if (!objectId || objectId === 'unknown' || objectId === 'demo') {
			loading = false;
			error =
				objectId === 'demo' ? $_('evidence.missing_demo') : $_('evidence.missing_id');
			meta = null;
			return;
		}
		loading = true;
		error = '';
		try {
			meta = await GetEvidence(objectId);
		} catch (e: any) {
			error = e?.message ?? String(e);
			meta = null;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		load(id);
	});

	function isVideo(ct: string) {
		return ct.startsWith('video/') || ct.includes('mp4') || ct.includes('webm');
	}
	function isAudio(ct: string) {
		return ct.startsWith('audio/') || ct.includes('wav') || ct.includes('ogg') || ct.includes('mpeg');
	}
	function isImage(ct: string) {
		return ct.startsWith('image/') || ct.includes('jpeg') || ct.includes('png') || ct.includes('jpg');
	}
	function formatBytes(n: number) {
		if (!n) return '\u2014';
		if (n < 1024) return `${n} B`;
		if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
		return `${(n / (1024 * 1024)).toFixed(1)} MB`;
	}
</script>

<button class="back-btn" onclick={() => push('/')}>{$_('common.back')}</button>

<div class="page-header">
	<h1 class="title"><span class="title-accent">{$_('evidence.title')}</span></h1>
	<p class="subtitle">
		object_id <code>{id || '\u2014'}</code>
		{#if meta}
			<span class="mode-chip" data-mode={meta.mode}>{meta.mode}</span>
		{/if}
	</p>
</div>

<div class="player-shell">
	{#if loading}
		<div class="player-placeholder"><span>{$_('common.loading')}</span></div>
	{:else if error}
		<div class="player-placeholder err-box">
			<span>{error}</span>
			<p class="hint">{$_('evidence.hint')}</p>
		</div>
	{:else if meta}
		{#if isVideo(meta.content_type)}
			<!-- svelte-ignore a11y_media_has_caption -->
			<video class="player" controls preload="metadata" src={meta.play_url}>
				{$_('evidence.no_video')}
			</video>
		{:else if isAudio(meta.content_type)}
			<audio class="player audio" controls preload="metadata" src={meta.play_url}>
				{$_('evidence.no_audio')}
			</audio>
		{:else if isImage(meta.content_type)}
			<img class="player img" src={meta.play_url} alt="Evidence {meta.object_id}" />
		{:else}
			<div class="player-placeholder">
				<span>{$_('evidence.unsupported', { type: meta.content_type || 'unknown' })}</span>
				<a class="dl" href={meta.play_url} target="_blank" rel="noopener">{$_('evidence.open_download')}</a>
			</div>
		{/if}

		<dl class="meta">
			<div><dt>{$_('evidence.kind')}</dt><dd>{meta.kind || '\u2014'}</dd></div>
			<div><dt>{$_('evidence.source')}</dt><dd>{meta.source || '\u2014'}</dd></div>
			<div><dt>{$_('evidence.type')}</dt><dd>{meta.content_type || '\u2014'}</dd></div>
			<div><dt>{$_('evidence.size')}</dt><dd>{formatBytes(meta.size_bytes)}</dd></div>
			<div><dt>{$_('evidence.captured')}</dt><dd>{meta.timestamp || '\u2014'}</dd></div>
			<div><dt>{$_('evidence.expires')}</dt><dd>{meta.expires_at || '\u2014'}</dd></div>
		</dl>

		<p class="note">
			{#if meta.mode === 'presign'}
				{$_('evidence.playback_presign', { mode: meta.mode })}
			{:else}
				{$_('evidence.playback_proxy', { mode: meta.mode })}
			{/if}
		</p>
	{/if}
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
		padding: 16px;
		text-align: center;
	}
	.err-box {
		color: #f87171;
	}
	.hint {
		font-size: 12px;
		color: var(--muted);
		line-height: 1.45;
		max-width: 420px;
		margin: 0;
	}
	.player {
		width: 100%;
		border-radius: var(--radius);
		border: 1px solid var(--glass-border);
		background: #000;
		--wails-draggable: no-drag;
	}
	.player.audio {
		min-height: 48px;
		padding: 12px;
		background: var(--glass);
	}
	.player.img {
		display: block;
		max-height: 420px;
		object-fit: contain;
		background: #0a0b12;
	}
	.mode-chip {
		margin-left: 8px;
		font-size: 10px;
		font-weight: 800;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		padding: 2px 8px;
		border-radius: 999px;
		border: 1px solid var(--glass-border);
		color: var(--muted);
	}
	.mode-chip[data-mode='presign'] {
		color: #4ade80;
		border-color: rgba(74, 222, 128, 0.4);
		background: rgba(74, 222, 128, 0.12);
	}
	.mode-chip[data-mode='proxy'] {
		color: #facc15;
		border-color: rgba(250, 204, 21, 0.35);
		background: rgba(250, 204, 21, 0.1);
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
		word-break: break-all;
	}
	.note {
		font-size: 12.5px;
		color: var(--muted);
		line-height: 1.5;
	}
	.dl {
		color: #ff7a8f;
		font-weight: 600;
	}
	code {
		font-size: 0.92em;
		background: rgba(255, 255, 255, 0.06);
		padding: 1px 5px;
		border-radius: 4px;
	}
</style>
