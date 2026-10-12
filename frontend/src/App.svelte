<script lang="ts">
	import './app.css';
	import { onMount } from 'svelte';
	import { Events, WML } from '@wailsio/runtime';
	import Router from 'svelte-spa-router';
	import { routes } from './lib/routes';

	const wailsVersion = 'v3.0.0-beta.5';
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
		<span class="brand-mark" data-wml-openURL="https://v3.wails.io">
			<img src="/wails.png" class="brand-logo" alt="Wails logo" />
		</span>
	</header>

	<!-- Routed content -->
	<Router {routes} />
</main>

<hr class="footer-divider" />

<footer class="footer">
	<span class="footer-version">{wailsVersion}</span>

	<span class="footer-time">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
			stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
			<circle cx="12" cy="12" r="10"></circle>
			<polyline points="12 6 12 12 16 14"></polyline>
		</svg>
		<span>{time}</span>
	</span>

	<a class="footer-docs" href="https://v3.wails.io" data-wml-openURL="https://v3.wails.io"
		aria-label="Wails documentation">
		Docs
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
			stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
			<line x1="7" y1="17" x2="17" y2="7"></line>
			<polyline points="7 7 17 7 17 17"></polyline>
		</svg>
	</a>
</footer>