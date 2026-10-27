<script lang="ts">
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { _ } from '../lib/i18n';
	import LangSwitcher from '../lib/LangSwitcher.svelte';

	/** Default local GEV dev/preview origin. Override with ?gev= or localStorage. */
	const DEFAULT_GEV = 'http://localhost:4173';

	let baseUrl = $state(DEFAULT_GEV);
	let embedSrc = $state('');
	let loadError = $state(false);
	let editing = $state(false);
	let draftUrl = $state(DEFAULT_GEV);

	function resolveBase(): string {
		try {
			const params = new URLSearchParams(window.location.hash.split('?')[1] || '');
			const fromQuery = params.get('gev');
			if (fromQuery) return fromQuery.replace(/\/$/, '');
		} catch {
			/* ignore */
		}
		try {
			const stored = localStorage.getItem('swail.gevBaseUrl');
			if (stored) return stored.replace(/\/$/, '');
		} catch {
			/* ignore */
		}
		return DEFAULT_GEV;
	}

	function buildEmbedSrc(base: string): string {
		const u = new URL(base.includes('://') ? base : `http://${base}`);
		u.searchParams.set('embed', '1');
		return u.toString();
	}

	function applyBase(next: string) {
		const cleaned = next.trim().replace(/\/$/, '') || DEFAULT_GEV;
		baseUrl = cleaned;
		embedSrc = buildEmbedSrc(cleaned);
		loadError = false;
		try {
			localStorage.setItem('swail.gevBaseUrl', cleaned);
		} catch {
			/* ignore */
		}
	}

	function openExternal() {
		const url = embedSrc || buildEmbedSrc(baseUrl);
		// Wails openURL when available; otherwise window.open
		const w = window as unknown as { wails?: { Browser?: { OpenURL?: (u: string) => void } } };
		if (w.wails?.Browser?.OpenURL) {
			w.wails.Browser.OpenURL(url);
		} else {
			window.open(url, '_blank', 'noopener,noreferrer');
		}
	}

	onMount(() => {
		applyBase(resolveBase());
		draftUrl = baseUrl;
	});
</script>

<header class="top-bar">
	<div class="top-left">
		<span class="brand-mark inline" data-wml-openURL="https://www.agentjetson.ai">
			<img src="/wails.png" class="brand-logo" alt="AgentJetson" />
		</span>
		<button class="back-btn" onclick={() => push('/')}>{$_('common.back')}</button>
		<div class="page-header compact">
			<h1 class="title"><span class="title-accent">{$_('map.title')}</span></h1>
			<p class="subtitle">{$_('map.subtitle')}</p>
		</div>
	</div>
	<div class="top-actions">
		<button class="ghost-btn" onclick={() => (editing = !editing)} type="button">
			{editing ? $_('map.hide_settings') : $_('map.settings')}
		</button>
		<button class="ghost-btn" onclick={openExternal} type="button">{$_('map.open_browser')}</button>
		<LangSwitcher />
	</div>
</header>

{#if editing}
	<div class="settings-bar">
		<label class="field">
			<span class="control-label">{$_('map.gev_url')}</span>
			<input
				class="input"
				type="url"
				bind:value={draftUrl}
				placeholder="http://localhost:4173"
				onkeydown={(e) => e.key === 'Enter' && applyBase(draftUrl)}
			/>
		</label>
		<button class="btn" type="button" onclick={() => applyBase(draftUrl)}>{$_('map.apply')}</button>
		<p class="hint">{$_('map.hint')}</p>
	</div>
{/if}

<div class="map-frame-wrap">
	{#if embedSrc && !loadError}
		<iframe
			class="map-frame"
			title="God's Eye View"
			src={embedSrc}
			allow="geolocation; microphone; clipboard-read; clipboard-write; fullscreen"
			referrerpolicy="no-referrer-when-downgrade"
			onerror={() => (loadError = true)}
		></iframe>
	{:else}
		<div class="map-fallback">
			<p class="empty">{$_('map.unavailable')}</p>
			<p class="hint">{$_('map.hint')}</p>
			<button class="btn" type="button" onclick={openExternal}>{$_('map.open_browser')}</button>
		</div>
	{/if}
</div>

<style>
	.top-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
		margin-bottom: var(--s-2);
		flex-shrink: 0;
	}

	.top-left {
		display: flex;
		align-items: center;
		gap: 0.65rem;
		min-width: 0;
	}

	.top-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.brand-mark.inline {
		display: inline-flex;
		align-items: center;
		flex-shrink: 0;
	}

	.brand-mark.inline .brand-logo {
		height: 36px;
		width: auto;
		display: block;
		object-fit: contain;
	}

	.page-header.compact {
		margin: 0;
	}

	.page-header.compact .title {
		margin: 0;
		font-size: clamp(1.35rem, 2.6vw, 1.7rem);
	}

	.page-header.compact .subtitle {
		margin: 2px 0 0;
		font-size: 13px;
	}

	.back-btn {
		margin-bottom: 0;
		flex-shrink: 0;
	}

	.ghost-btn {
		background: transparent;
		border: 1px solid var(--glass-border);
		color: var(--muted);
		padding: 6px 12px;
		border-radius: 8px;
		font-size: 13px;
		font-weight: 600;
		cursor: pointer;
		--wails-draggable: no-drag;
		transition: color 0.2s ease, border-color 0.2s ease;
	}

	.ghost-btn:hover {
		color: var(--text);
		border-color: rgba(255, 77, 77, 0.45);
	}

	.settings-bar {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		gap: 10px;
		margin-bottom: var(--s-2);
		padding: 12px 14px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
		--wails-draggable: no-drag;
	}

	.field {
		flex: 1 1 220px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.control-label {
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--muted);
	}

	.input {
		background: rgba(0, 0, 0, 0.25);
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

	.hint {
		flex: 1 1 100%;
		margin: 0;
		font-size: 12px;
		color: var(--muted);
		line-height: 1.45;
	}

	.map-frame-wrap {
		flex: 1 1 auto;
		min-height: 0;
		width: 100%;
		border-radius: var(--radius);
		border: 1px solid var(--glass-border);
		overflow: hidden;
		background: rgba(0, 0, 0, 0.35);
		position: relative;
	}

	.map-frame {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		border: 0;
		background: #0a0c14;
	}

	.map-fallback {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: center;
		gap: 12px;
		padding: var(--s-4);
		height: 100%;
		min-height: 280px;
	}

	.empty {
		margin: 0;
		color: var(--muted);
		font-size: 14px;
	}

	@media (max-width: 640px) {
		.top-bar {
			flex-wrap: wrap;
		}
		.top-actions {
			width: 100%;
			justify-content: flex-end;
		}
	}
</style>
