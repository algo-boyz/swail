<script lang="ts">
	import { push } from 'svelte-spa-router';

	function selectPlatform(platform: string) {
		push(`/${platform}`);
	}
</script>

<h1 class="title">
	<span class="title-accent">Swail</span>
	<span class="title-name">
		<span class="title-name-text">SniGo Dashboard</span>
	</span>
</h1>

<p class="subtitle">Cross platform arbitrage sniper app.</p>

<div class="platforms">
	<button
		class="platform-card"
		onclick={() => selectPlatform('marktplaats')}
		aria-label="Select Marktplaats"
	>
		<span class="platform-name">Marktplaats</span>
		<span class="platform-arrow">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"
				stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
				<line x1="5" y1="12" x2="19" y2="12"></line>
				<polyline points="12 5 19 12 12 19"></polyline>
			</svg>
		</span>
	</button>

	<button
		class="platform-card"
		onclick={() => selectPlatform('catawiki')}
		aria-label="Select Catawiki"
	>
		<span class="platform-name">Catawiki</span>
		<span class="platform-arrow">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"
				stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
				<line x1="5" y1="12" x2="19" y2="12"></line>
				<polyline points="12 5 19 12 12 19"></polyline>
			</svg>
		</span>
	</button>
</div>