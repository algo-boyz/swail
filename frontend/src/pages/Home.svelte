<script lang="ts">
	import { push } from 'svelte-spa-router';

	const cards = [
		{ path: '/query', title: 'Query', desc: 'Plates, objects, scenes, transcripts' },
		{ path: '/cameras', title: 'Cameras', desc: 'Sources & live tiles' },
		{ path: '/alerts', title: 'Alerts', desc: 'Recent detections & intents' },
		{ path: '/evidence/demo', title: 'Evidence', desc: 'Presigned audio / video player' }
	];
</script>

<div class="home">
	<span class="badge">scaffold · PR-1</span>

	<h1 class="title">
		<span class="title-accent">AgentJetson</span>
	</h1>
	<p class="subtitle">
		Intelligence at the edge of every encounter — voice query, evidence & awareness.
	</p>

	<nav class="platforms" aria-label="Main sections">
		{#each cards as c}
			<button class="platform-card" onclick={() => push(c.path)} aria-label={c.title}>
				<span>
					<span class="platform-name">{c.title}</span>
					<span class="card-desc">{c.desc}</span>
				</span>
				<span class="platform-arrow">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"
						stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
						<line x1="5" y1="12" x2="19" y2="12"></line>
						<polyline points="12 5 19 12 12 19"></polyline>
					</svg>
				</span>
			</button>
		{/each}
	</nav>
</div>

<style>
	.home {
		width: 100%;
	}
	.card-desc {
		display: block;
		font-size: 12.5px;
		font-weight: 500;
		color: var(--muted);
		margin-top: 3px;
	}
	.platform-card {
		align-items: flex-start;
	}
</style>
