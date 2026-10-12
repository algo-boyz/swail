<script lang="ts">
	import './app.css';
	import { onMount } from 'svelte';
	import { Events, WML } from '@wailsio/runtime';
	import Router from 'svelte-spa-router';
	import { routes } from './lib/routes';

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
</style>
