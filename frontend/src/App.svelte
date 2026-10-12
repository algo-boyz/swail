<script lang="ts">
	import './app.css';
	import { onMount } from 'svelte';
	import { Events, WML } from '@wailsio/runtime';
	import Router from 'svelte-spa-router';
	import { routes } from './lib/routes';
	import { locale, locales, setLocale, _ } from './lib/i18n';
	import type { Locale } from './lib/i18n';

	const appVersion = 'scaffold-pr1';
	let time: string = $state('Listening for Time event...');

	onMount(() => {
		Events.On('time', (v: any) => {
			const full = v.data;
			const compact = (full.match(/\d{1,2}:\d{2}:\d{2}/) || [full])[0];
			time = window.matchMedia('(max-width: 640px)').matches ? compact : full;
		});
		WML.Reload();
	});

	function onLangChange(e: Event) {
		const code = (e.currentTarget as HTMLSelectElement).value as Locale;
		setLocale(code);
	}
</script>

<main class="container">
	<header class="brand">
		<span class="brand-mark" data-wml-openURL="https://www.agentjetson.ai">
			<img src="/wails.png" class="brand-logo" alt="AgentJetson" />
		</span>
	</header>

	<div class="route-slot">
		<Router {routes} />
	</div>
</main>

<hr class="footer-divider" />

<footer class="footer">
	<span class="footer-version">{appVersion}</span>

	<label class="lang-switch" title={$_('lang.label')}>
		<span class="lang-label">{$_('lang.label')}</span>
		<select class="lang-select" value={$locale} onchange={onLangChange} aria-label={$_('lang.label')}>
			{#each locales as l}
				<option value={l.code}>{l.label}</option>
			{/each}
		</select>
	</label>

	<span class="footer-time">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
			stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
			<circle cx="12" cy="12" r="10"></circle>
			<polyline points="12 6 12 12 16 14"></polyline>
		</svg>
		<span>{time}</span>
	</span>

	<a
		class="footer-docs"
		href="https://www.agentjetson.ai"
		data-wml-openURL="https://www.agentjetson.ai"
		aria-label="AgentJetson site"
	>
		agentjetson.ai
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
			stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
			<line x1="7" y1="17" x2="17" y2="7"></line>
			<polyline points="7 7 17 7 17 17"></polyline>
		</svg>
	</a>
</footer>

<style>
	.route-slot {
		flex: 1 1 auto;
		min-height: 0;
		width: 100%;
		display: flex;
		flex-direction: column;
	}

	/* svelte-spa-router mounts the page as a direct child; stretch it */
	.route-slot :global(> *) {
		flex: 1 1 auto;
		min-height: 0;
		width: 100%;
		display: flex;
		flex-direction: column;
	}

	.lang-switch {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		font-size: 12px;
		color: var(--muted);
		--wails-draggable: no-drag;
	}

	.lang-label {
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		font-size: 10px;
	}

	.lang-select {
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 8px;
		color: var(--text);
		padding: 4px 8px;
		font-size: 12.5px;
		font-weight: 600;
		cursor: pointer;
		outline: none;
	}

	.lang-select:focus {
		border-color: rgba(255, 77, 77, 0.55);
	}
</style>
